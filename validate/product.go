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
	AddAlias("xProductName", "min=1,max=50")
	AddAlias("xProductPrice", "min=0.01,max=1000")
	AddAlias("xProductUnit", "min=1,max=20")
	AddAlias("xProductCatalog", "min=1,max=1000")
	AddAlias("xProductSN", "min=1,max=100")
	AddAlias("xProductMainPic", "min=1,max=20")
	AddAlias("xProductKeyword", "min=1,max=100")
	AddAlias("xProductOrigin", "min=1,max=100")
	AddAlias("xProductBrand", "min=1")

	AddAlias("xProductCategoryName", "min=1,max=10")
	AddAlias("xProductCategoryLevel", "number,min=1,max=3")
}
