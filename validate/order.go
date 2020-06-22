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

import "github.com/go-playground/validator/v10"

func init() {
	// 暂时订单状态仅设置为最大10
	AddAlias("xOrderStatus", "number,min=1,max=10")
	// 暂时子订单状态仅设置为最大10
	AddAlias("xSubOrderStatus", "number,min=1,max=10")
	AddAlias("xOrderProductID", "number,min=1")
	// 暂仅支持一次最大1000
	AddAlias("xOrderProductCount", "number,min=1,max=1000")
	// 订单编号
	AddAlias("xOrderSN", "min=1")
	// 订单金额
	AddAlias("xOrderAmount", "number")
	// 订单客户id
	AddAlias("xOrderUser", "number")
	// 订单运输单编号
	AddAlias("xOrderDeliverySN", "min=10,max=20")
	// 订单运输公司
	AddAlias("xOrderDeliveryCompnay", "min=1,max=10")
	// 订单送货人
	AddAlias("xOrderCourier", "number,min=1")

	// 支付来源
	Add("xSource", func(fl validator.FieldLevel) bool {
		return isInString(fl, []string{
			"wechat",
			"alipay",
		})
	})
}
