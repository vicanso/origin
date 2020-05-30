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

package cs

import "strconv"

const (
	// CID context id
	CID = "cid"
	// UserSession user session
	UserSession = "userSession"
)

const (
	// ConfigEnabled config enabled
	ConfigEnabled = iota + 1
	// ConfigDiabled config disabled
	ConfigDiabled
)

const (
	// MagicalCaptcha magical captcha(for test only)
	MagicalCaptcha = "0145"
)

// 用户状态
const (
	// AccountStatusEnabled account enabled
	AccountStatusEnabled = iota + 1
	// AccountStatusForbidden account forbidden
	AccountStatusForbidden
)

var (
	AccountStatuses = []int{
		AccountStatusEnabled,
		AccountStatusForbidden,
	}
	AccountStatusesString = []string{
		strconv.Itoa(AccountStatusEnabled),
		strconv.Itoa(AccountStatusForbidden),
	}
)

// 用户角色
const (
	// UserRoleNormal normal user
	UserRoleNormal = "normal"
	// UserRoleSu super user
	UserRoleSu = "su"
	// UserRoleAdmin admin user
	UserRoleAdmin = "admin"
)

var (
	UserRoles = []string{
		UserRoleNormal,
		UserRoleSu,
		UserRoleAdmin,
	}
)

// 用户群组
const (
	UserGroupIT        = "it"
	UserGroupFinance   = "finance"
	UserGroupMarketing = "marketing"
)

// 用户分组
var (
	UserGroups = []string{
		UserGroupIT,
		UserGroupFinance,
		UserGroupMarketing,
	}
)

// 品牌状态
const (
	BrandStatusEnabled = iota + 1
	BrandStatusDisabled
)

// 品牌状态列表
var (
	BrandStatuses = []int{
		BrandStatusEnabled,
		BrandStatusDisabled,
	}
)

// 产品状态
const (
	ProductStatusEnabled = iota + 1
	ProductStatusDisabled
)

// 产品状态列表
var (
	ProductStatuses = []int{
		ProductStatusEnabled,
		ProductStatusDisabled,
	}
)

// 地区状态
const (
	RegionStatusEnabled = iota + 1
	RegionStatusDisabled
)

var (
	RegionStatuses = []int{
		RegionStatusEnabled,
		RegionStatusDisabled,
	}
	RegionStatusesString = []string{
		strconv.Itoa(RegionStatusEnabled),
		strconv.Itoa(RegionStatusDisabled),
	}
)

// 地区分类
const (
	RegionCountry  = "country"
	RegionProvince = "province"
	RegionCity     = "city"
	RegionArea     = "area"
	RegionStreet   = "street"
)

var (
	RegionCategories = []string{
		RegionCountry,
		RegionProvince,
		RegionCity,
		RegionArea,
		RegionStreet,
	}
)
