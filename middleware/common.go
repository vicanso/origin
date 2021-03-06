// Copyright 2019 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package middleware

import (
	"net/http"
	"strings"
	"time"

	warner "github.com/vicanso/count-warner"
	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	"github.com/vicanso/origin/helper"
	"github.com/vicanso/origin/service"
	"github.com/vicanso/tiny/log"
	"go.uber.org/zap"
)

const (
	xCaptchHeader     = "X-Captcha"
	errCommonCategory = "common-validate"
)

var (
	errQueryNotAllow = &hes.Error{
		StatusCode: http.StatusBadRequest,
		Message:    "query is not allowed",
		Category:   errCommonCategory,
	}
	errCaptchaIsInvalid = &hes.Error{
		StatusCode: http.StatusBadRequest,
		Message:    "图形验证码错误",
		Category:   errCommonCategory,
	}
	errCaptchaExpired = &hes.Error{
		StatusCode: http.StatusBadRequest,
		Message:    "图形验证码已过期，请刷新",
		Category:   errCommonCategory,
	}
)

// NoQuery no query middleware
func NoQuery(c *elton.Context) (err error) {
	if c.Request.URL.RawQuery != "" {
		err = errQueryNotAllow
		return
	}
	return c.Next()
}

// WaitFor at least wait for duration
func WaitFor(d time.Duration, args ...bool) elton.Handler {
	ns := d.Nanoseconds()
	onlyErrOccurred := false
	if len(args) != 0 {
		onlyErrOccurred = args[0]
	}
	return func(c *elton.Context) (err error) {
		start := time.Now()
		err = c.Next()
		// 如果未出错，而且配置为仅在出错时才等待
		if err == nil && onlyErrOccurred {
			return
		}
		use := time.Now().UnixNano() - start.UnixNano()
		// 无论成功还是失败都wait for
		if use < ns {
			time.Sleep(time.Duration(ns-use) * time.Nanosecond)
		}
		return
	}
}

// ValidateCaptcha validate chapter
func ValidateCaptcha(magicalCaptcha string) elton.Handler {
	return func(c *elton.Context) (err error) {
		value := c.GetRequestHeader(xCaptchHeader)
		if value == "" {
			err = errCaptchaIsInvalid
			return
		}
		arr := strings.Split(value, ":")
		if len(arr) != 2 {
			err = errCaptchaIsInvalid
			return
		}
		// 如果有配置万能验证码，则判断是否相等
		if magicalCaptcha != "" && arr[1] == magicalCaptcha {
			return c.Next()
		}
		valid, err := service.ValidateCaptcha(arr[0], arr[1])
		if err != nil {
			if helper.IsRedisNilError(err) {
				err = errCaptchaExpired
			}
			return err
		}
		if !valid {
			err = errCaptchaIsInvalid
			return
		}
		return c.Next()
	}
}

// NewNoCacheWithCondition create a nocache middleware
func NewNoCacheWithCondition(key, value string) elton.Handler {
	return func(c *elton.Context) (err error) {
		err = c.Next()
		if c.QueryParam(key) == value {
			c.NoCache()
		}
		return
	}
}

// NewNotFoundHandler new not found handler
func NewNotFoundHandler() http.HandlerFunc {
	// 对于404的请求，不会执行中间件，一般都是因为攻击之类才会导致大量出现404，
	// 因此可在此处汇总出错IP，针对较频繁出错IP，增加告警信息
	// 如果1分钟同一个IP出现60次404
	warner404 := warner.NewWarner(60*time.Second, 60)
	warner404.ResetOnWarn = true
	warner404.On(func(ip string, _ warner.Count) {
		service.AlarmError("too many 404 request, client ip:" + ip)
	})
	go func() {
		// 因为404是根据IP来告警，因此可能存在大量不同的key，因此定时清除过期数据
		for range time.NewTicker(5 * time.Minute).C {
			warner404.ClearExpired()
		}
	}()
	logger := log.Default()
	notFoundErr := &hes.Error{
		Message:    "Not Found",
		StatusCode: http.StatusNotFound,
		Category:   "defaultNotFound",
	}
	notFoundErrBytes := notFoundErr.ToJSON()
	return func(resp http.ResponseWriter, req *http.Request) {
		ip := elton.GetClientIP(req)
		logger.Info("404",
			zap.String("ip", ip),
			zap.String("method", req.Method),
			zap.String("uri", req.RequestURI),
		)
		resp.Header().Set(elton.HeaderContentType, elton.MIMEApplicationJSON)
		resp.WriteHeader(http.StatusNotFound)
		_, err := resp.Write(notFoundErrBytes)
		if err != nil {
			logger.Info("404 response fail",
				zap.String("ip", ip),
				zap.String("uri", req.RequestURI),
				zap.Error(err),
			)
		}
		warner404.Inc(ip, 1)
	}
}

// NewMethodNotAllowedHandler new method not allow handler
func NewMethodNotAllowedHandler() http.HandlerFunc {
	logger := log.Default()
	methodNotAllowedErr := &hes.Error{
		Message:    "Method Not Allowed",
		StatusCode: http.StatusMethodNotAllowed,
		Category:   "defaultMethodNotAllowed",
	}
	methodNotAllowedErrBytes := methodNotAllowedErr.ToJSON()
	return func(resp http.ResponseWriter, req *http.Request) {
		ip := elton.GetClientIP(req)
		logger.Info("method not allowed",
			zap.String("ip", ip),
			zap.String("method", req.Method),
			zap.String("uri", req.RequestURI),
		)
		resp.Header().Set(elton.HeaderContentType, elton.MIMEApplicationJSON)
		resp.WriteHeader(http.StatusMethodNotAllowed)
		_, err := resp.Write(methodNotAllowedErrBytes)
		if err != nil {
			logger.Info("method not allowed response fail",
				zap.String("ip", ip),
				zap.String("uri", req.RequestURI),
				zap.Error(err),
			)
		}
	}
}
