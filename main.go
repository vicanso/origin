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

/*
Package main Origin Server

	This should demonstrate all the possible comment annotations
	that are available to turn go code into a fully compliant swagger 2.0 spec

Host: localhost
BasePath: /
Version: 1.0.0
Schemes: http

Consumes:
- application/json

Produces:
- application/json

swagger:meta
*/
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	warner "github.com/vicanso/count-warner"
	"github.com/vicanso/elton"
	compress "github.com/vicanso/elton-compress"
	M "github.com/vicanso/elton/middleware"
	"github.com/vicanso/hes"
	"github.com/vicanso/origin/config"
	_ "github.com/vicanso/origin/controller"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/helper"
	"github.com/vicanso/origin/log"
	"github.com/vicanso/origin/middleware"
	"github.com/vicanso/origin/router"
	_ "github.com/vicanso/origin/schedule"
	"github.com/vicanso/origin/service"
	"github.com/vicanso/origin/util"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	// Version version of tiny
	Version string
	// BuildAt build at
	BuildedAt string
)

func init() {
	_, _ = maxprocs.Set(maxprocs.Logger(func(format string, args ...interface{}) {
		value := fmt.Sprintf(format, args...)
		log.Default().Info(value)
	}))
	config.SetVersion(Version)
	config.SetBuildedAt(BuildedAt)
}

// 相关依赖服务的校验，主要是数据库等
func dependServiceCheck() (err error) {
	err = helper.RedisPing()
	if err != nil {
		return
	}
	configSrv := new(service.ConfigurationSrv)
	err = configSrv.Refresh()
	if err != nil {
		return
	}
	return
}

func main() {
	closeOnce := sync.Once{}
	closeDeps := func() {
		// 关闭influxdb，flush统计数据
		helper.GetInfluxSrv().Close()
	}
	defer func() {
		closeOnce.Do(closeDeps)
	}()

	logger := log.Default()
	e := elton.New()

	// 启用耗时跟踪
	if util.IsDevelopment() {
		e.EnableTrace = true
	}
	e.OnTrace(func(c *elton.Context, infos elton.TraceInfos) {
		// 设置server timing
		c.ServerTiming(infos, "origni-")
	})

	// 是否用户关闭
	closedByUser := false
	// 非开发环境，监听信号退出
	if !util.IsDevelopment() {

		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			for s := range c {
				logger.Info("server will be closed",
					zap.String("signal", s.String()),
				)
				closedByUser = true
				// 设置状态为退出中，/ping请求返回出错，反向代理不再转发流量
				config.SetApplicationStatus(config.ApplicationStatusStopping)
				// docker 在10秒内退出，因此设置8秒后退出
				time.Sleep(5 * time.Second)
				// 所有新的请求均返回出错
				e.GracefulClose(3 * time.Second)
				closeOnce.Do(closeDeps)
				os.Exit(0)
			}
		}()
	}

	e.SignedKeys = service.GetSignedKeys()
	e.GenerateID = func() string {
		return util.RandomString(8)
	}

	// 未处理的error才会触发
	// 如果1分钟出现超过5次未处理异常
	// exception的warner只有一个key，因此无需定时清除
	warnerException := warner.NewWarner(60*time.Second, 5)
	warnerException.ResetOnWarn = true
	warnerException.On(func(_ string, _ warner.Count) {
		service.AlarmError("too many uncaught exception")
	})
	e.OnError(func(c *elton.Context, err error) {
		he := hes.Wrap(err)
		if !util.IsProduction() {
			if he.Extra == nil {
				he.Extra = make(map[string]interface{})
			}
			he.Extra["stack"] = util.GetStack(5)
		}
		ip := c.RealIP()
		uri := c.Request.RequestURI

		helper.GetInfluxSrv().Write(cs.MeasurementException, map[string]interface{}{
			"ip":  ip,
			"uri": uri,
		}, map[string]string{
			"category": "routeError",
		})

		// 可以针对实际场景输出更多的日志信息
		logger.DPanic("exception",
			zap.String("ip", ip),
			zap.String("uri", uri),
			zap.Error(he.Err),
		)
		warnerException.Inc("exception", 1)
	})

	e.NotFoundHandler = middleware.NewNotFoundHandler()
	e.MethodNotAllowedHandler = middleware.NewMethodNotAllowedHandler()

	// 捕捉panic异常，避免程序崩溃
	e.UseWithName(M.NewRecover(), "recover")

	e.UseWithName(middleware.NewEntry(), "entry")

	// 接口相关统计信息
	e.UseWithName(middleware.NewStats(), "stats")

	// 限制最大请求量
	maxRequestLimit := config.GetRequestLimit()
	if maxRequestLimit != 0 {
		e.UseWithName(M.NewGlobalConcurrentLimiter(M.GlobalConcurrentLimiterConfig{
			Max: maxRequestLimit,
		}), "request-limit")
	}

	// 配置只针对snappy与lz4压缩（主要用于减少内网线路带宽，对外的压缩由前置反向代理 完成）
	compressMinLength := 2 * 1024
	compressConfig := M.NewCompressConfig(
		&compress.SnappyCompressor{
			MinLength: compressMinLength,
		},
		&compress.Lz4Compressor{
			MinLength: compressMinLength,
		},
	)
	e.UseWithName(M.NewCompress(compressConfig), "compress")

	// 错误处理，将错误转换为json响应
	e.UseWithName(middleware.NewError(), "error")

	// IP限制
	e.UseWithName(middleware.NewIPBlocker(), "ip-blocker")

	// 根据配置对路由mock返回
	e.UseWithName(middleware.NewRouterMocker(), "router-mocker")

	// 路由并发限制
	e.UseWithName(M.NewRCL(M.RCLConfig{
		Limiter: service.GetRouterConcurrencyLimiter(),
	}), "rcl")

	// etag与fresh的处理
	e.UseWithName(M.NewDefaultFresh(), "fresh").
		UseWithName(M.NewDefaultETag(), "etag")

	// 对响应数据 c.Body 转换为相应的json响应
	e.UseWithName(M.NewDefaultResponder(), "responder")

	// 读取读取body的数的，转换为json bytes
	e.UseWithName(M.NewDefaultBodyParser(), "body-parser")

	// 初始化路由
	for _, g := range router.GetGroups() {
		e.AddGroup(g)
	}

	service.InitRouterConcurrencyLimiter(e.Routers)

	err := dependServiceCheck()
	if err != nil {
		service.AlarmError("check depend service fail, " + err.Error())
		// 可以针对实际场景输出更多的日志信息
		logger.DPanic("exception",
			zap.Error(err),
		)
		panic(err)
	}

	listen := config.GetListen()
	// http1与http2均支持
	e.Server = &http.Server{
		Handler: h2c.NewHandler(e, &http2.Server{}),
	}

	logger.Info("server will listen on " + listen)
	config.SetApplicationStatus(config.ApplicationStatusRunning)
	err = e.ListenAndServe(listen)

	// 如果出错而且非用户关闭，则发送告警
	if err != nil && !closedByUser {
		service.AlarmError("listen and serve fail, " + err.Error())
		panic(err)
	}
}
