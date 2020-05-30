// Copyright 2020 tree xie
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

func init() {
	AddAlias("xBrandName", "min=1,max=20")
	AddAlias("xBrandLogo", "min=1,max=100")
	AddAlias("xBrandCatalog", "min=1,max=1000")

	// brandStatusesStr := make([]string, 0)
	// for _, item := range cs.BrandStatuses {
	// 	brandStatusesStr = append(brandStatusesStr, strconv.Itoa(item))
	// }

	// Add("xStatus", func(fl validator.FieldLevel) bool {
	// 	if isInt(fl) {
	// 		return isInInt(fl, cs.BrandStatuses)
	// 	}
	// 	return isInString(fl, brandStatusesStr)
	// })
}
