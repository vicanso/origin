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

package validate

import (
	"github.com/go-playground/validator/v10"

	"github.com/vicanso/origin/cs"
)

func init() {
	// 账号
	AddAlias("xUserAccount", "ascii,min=2,max=10")

	AddAlias("xUserPassword", "ascii,len=44")

	AddAlias("xUserName", "min=1,max=20")

	AddAlias("xUserID", "numeric")

	AddAlias("xUserMarketingGroup", "ascii,min=1,max=10")

	Add("xUserRole", func(fl validator.FieldLevel) bool {
		return isInString(fl, cs.UserRoles)
	})
	Add("xUserRoles", func(fl validator.FieldLevel) bool {
		return isAllInString(fl, cs.UserRoles)
	})
	Add("xUserGroup", func(fl validator.FieldLevel) bool {
		return isInString(fl, cs.UserGroups)
	})
	Add("xUserGroups", func(fl validator.FieldLevel) bool {
		return isAllInString(fl, cs.UserGroups)
	})
	AddAlias("xUserEmail", "email")

	AddAlias("xUserTrackCategory", "min=1,max=20")
}
