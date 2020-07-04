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
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/vicanso/origin/cs"
)

func init() {
	// 字符串，max表示字符串长度
	AddAlias("xLimit", "number,min=1,max=2")
	AddAlias("xOffset", "number,min=0,max=5")
	AddAlias("xOrder", "ascii,min=0,max=100")
	AddAlias("xFields", "ascii,min=0,max=100")
	AddAlias("xKeyword", "min=1,max=10")
	AddAlias("xMobile", "number,len=11")
	// 地址编码
	AddAlias("xBaseAddress", "min=1,max=10")
	// 详细地址
	AddAlias("xAddress", "min=1,max=100")
	// 优先级
	AddAlias("xPriority", "min=1,max=100")
	// 排序
	AddAlias("xRank", "min=1,max=1000")
	// 文件路径，以/files为前缀的uri
	AddAlias("xFile", "startswith=/files")

	durationRegexp := regexp.MustCompile("^[1-9][0-9]*(ms|[smh])$")
	Add("xDuration", func(fl validator.FieldLevel) bool {
		value, ok := toString(fl)
		if !ok {
			return false
		}
		return durationRegexp.MatchString(value)
	})

	Add("xStatus", func(fl validator.FieldLevel) bool {
		if isInt(fl) {
			return isInInt(fl, cs.Statuses)
		}
		return isInString(fl, cs.StatusesString)
	})
}
