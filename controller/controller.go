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

package controller

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/helper"
	"github.com/vicanso/origin/log"
	"github.com/vicanso/origin/middleware"
	"github.com/vicanso/origin/service"
	"github.com/vicanso/origin/util"

	"go.uber.org/zap"

	M "github.com/vicanso/elton/middleware"
)

var (
	errShouldLogin  = hes.New("should login first")
	errLoginAlready = hes.New("login already, please logout first")
	errForbidden    = &hes.Error{
		StatusCode: http.StatusForbidden,
		Message:    "禁止使用该功能",
	}
)

var (
	logger     = log.Default()
	now        = util.NowString
	getTrackID = util.GetTrackID

	// 服务列表
	// 配置服务
	configSrv = new(service.ConfigurationSrv)
	// 用户服务
	userSrv = new(service.UserSrv)
	// 品牌服务
	brandSrv = new(service.BrandSrv)
	// 文件服务
	fileSrv = new(service.FileSrv)
	// 产品服务
	productSrv = new(service.ProductSrv)
	// 地区服务
	regionSrv = new(service.RegionSrv)
	// 订单服务
	orderSrv = new(service.OrderSrv)
	// 供应商服务
	supplierSrv = new(service.SupplierSrv)
	// 收货人服务
	receiverSrv = new(service.ReceiverSrv)
	// 广告服务
	advertisementSrv = new(service.AdvertisementSrv)
	// 图片服务
	imageSrv = new(service.ImageSrv)

	// 创建新的并发控制中间件
	newConcurrentLimit = middleware.NewConcurrentLimit
	// 创建IP限制中间件
	newIPLimit = middleware.NewIPLimit
	// 创建出错限制中间件
	newErrorLimit = middleware.NewErrorLimit

	getUserSession = service.NewUserSession
	// 加载用户session
	loadUserSession = middleware.NewSession()
	// 判断用户是否登录
	shouldBeLogined = checkLogin
	// 判断用户是否未登录
	shouldBeAnonymous = checkAnonymous
	// 判断用户是否admin权限
	shouldBeAdmin = newCheckRolesMiddleware([]string{
		cs.UserRoleSu,
		cs.UserRoleAdmin,
	})
	// shouldBeSu 判断用户是否su权限
	shouldBeSu = newCheckRolesMiddleware([]string{
		cs.UserRoleSu,
	})
	// noCacheIfSetNoCache 如果query指定了no cache，则设置不缓存
	noCacheIfSetNoCache = middleware.NewNoCacheWithCondition("cacheControl", "no-cache")

	// checkMarketingGroup 判断用户组是否marketing
	checkMarketingGroup = newCheckGroupsMiddleware([]string{
		cs.UserGroupMarketing,
	})
	// checkLogisticsGroup 判断用户组是否logistics
	checkLogisticsGroup = newCheckGroupsMiddleware([]string{
		cs.UserGroupLogistics,
	})

	// 图形验证码校验
	captchaValidate elton.Handler
)

type (
	// listParams list params
	listParams struct {
		Limit  string `json:"limit,omitempty" validate:"xLimit"`
		Offset string `json:"offset,omitempty" validate:"omitempty,xOffset"`
		Fields string `json:"fields,omitempty" validate:"omitempty,xFields"`
		Order  string `json:"order,omitempty" validate:"omitempty,xOrder"`
	}
	// queryConditions query conditions
	queryConditions struct {
		queryList []string
		args      []interface{}
	}
)

type PGQueryParams = helper.PGQueryParams

func init() {
	magicalValue := ""
	if !util.IsProduction() {
		magicalValue = cs.MagicalCaptcha
	}
	captchaValidate = middleware.ValidateCaptcha(magicalValue)
}

// add add query and argument to query condition
func (conditions *queryConditions) add(query string, arg interface{}) {
	conditions.addQuery(query)
	conditions.addArgs(arg)
}

// addQuery add query to query condition
func (conditions *queryConditions) addQuery(arr ...string) {
	if len(conditions.queryList) == 0 {
		conditions.queryList = make([]string, 0)
	}
	conditions.queryList = append(conditions.queryList, arr...)
}

// addArgs add arguments to query condition
func (conditions *queryConditions) addArgs(args ...interface{}) {
	if len(conditions.args) == 0 {
		conditions.args = make([]interface{}, 0)
	}
	conditions.args = append(conditions.args, args...)
}

// toArray convert query conditions to gorm query arguments
func (conditions *queryConditions) toArray() []interface{} {
	arr := make([]interface{}, 0)
	if len(conditions.queryList) != 0 {
		arr = append(arr, strings.Join(conditions.queryList, " AND "))
		arr = append(arr, conditions.args...)
	}
	return arr
}

// newTracker create a new trakcer middleware
func newTracker(action string) elton.Handler {
	return M.NewTracker(M.TrackerConfig{
		Mask: regexp.MustCompile(`(?i)password`),
		OnTrack: func(info *M.TrackerInfo, c *elton.Context) {
			account := ""
			us := service.NewUserSession(c)
			if us != nil && us.IsLogined() {
				account = us.GetAccount()
			}
			fields := make([]zap.Field, 0, 10)
			fields = append(
				fields,
				zap.String("action", action),
				zap.String("cid", info.CID),
				zap.String("account", account),
				zap.String("ip", c.RealIP()),
				zap.String("sid", util.GetSessionID(c)),
				zap.Int("result", info.Result),
			)
			if info.Query != nil {
				fields = append(fields, zap.Any("query", info.Query))
			}
			if info.Params != nil {
				fields = append(fields, zap.Any("params", info.Params))
			}
			if info.Form != nil {
				fields = append(fields, zap.Any("form", info.Form))
			}
			if info.Err != nil {
				fields = append(fields, zap.Error(info.Err))
			}
			logger.Info("tracker", fields...)
		},
	})
}

func isLogined(c *elton.Context) bool {
	us := service.NewUserSession(c)
	return us.IsLogined()
}

func checkLogin(c *elton.Context) (err error) {
	if !isLogined(c) {
		err = errShouldLogin
		return
	}
	return c.Next()
}

func checkAnonymous(c *elton.Context) (err error) {
	if isLogined(c) {
		err = errLoginAlready
		return
	}
	return c.Next()
}

// newCheckRolesMiddleware create a new check roles middleware
func newCheckRolesMiddleware(validRoles []string) elton.Handler {
	return func(c *elton.Context) (err error) {
		if !isLogined(c) {
			err = errShouldLogin
			return
		}
		us := service.NewUserSession(c)
		roles := us.GetRoles()
		valid := util.UserRoleIsValid(validRoles, roles)
		if valid {
			return c.Next()
		}
		err = errForbidden
		return
	}
}

// newCheckGroupsMiddleware create a new check groups middleware
func newCheckGroupsMiddleware(validGroups []string) elton.Handler {
	return func(c *elton.Context) (err error) {
		if !isLogined(c) {
			err = errShouldLogin
			return
		}
		us := service.NewUserSession(c)
		groups := us.GetGroups()
		valid := util.UserGroupIsValid(validGroups, groups)
		if valid {
			return c.Next()
		}
		err = errForbidden
		return
	}
}

// toPGQueryParams to pg query params
func (params listParams) toPGQueryParams() PGQueryParams {
	limit, _ := strconv.Atoi(params.Limit)
	offset, _ := strconv.Atoi(params.Offset)
	return PGQueryParams{
		Limit:  limit,
		Offset: offset,
		Order:  params.Order,
		Fields: params.Fields,
	}
}

// getIDFromParams get id form context params
func getIDFromParams(c *elton.Context) (id uint, err error) {
	value, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		he := hes.Wrap(err)
		he.Category = "parse-int"
		err = he
		return
	}
	id = uint(value)
	return
}
