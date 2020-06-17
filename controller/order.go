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

package controller

import (
	"net/http"
	"time"

	"github.com/vicanso/elton"
	M "github.com/vicanso/elton/middleware"
	"github.com/vicanso/hes"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/middleware"
	"github.com/vicanso/origin/router"
	"github.com/vicanso/origin/service"
	"github.com/vicanso/origin/validate"
)

type (
	orderCtrl struct{}

	addOrderParams struct {
		Products []struct {
			ID    uint `json:"id,omitempty" validate:"xOrderProductID"`
			Count uint `json:"count,omitempty" validate:"xOrderProductCount"`
		} `json:"products,omitempty"`
	}
	// 支付金额
	payOrderParams struct {
		PayAmount float64 `json:"payAmount,omitempty" validate:"xOrderAmount"`
		Source    string  `json:"source,omitempty" validate:"xSource"`
	}
	listOrderParams struct {
		listParams

		Status string `json:"status,omitempty" validate:"omitempty,xOrderStatus"`
		SN     string `json:"sn,omitempty" validate:"omitempty,xOrderSN"`
	}
)

const (
	errOrderCtrlCategory = "order-ctrl"
)

var (
	errProductsIsEmpty = &hes.Error{
		Message:    "产品不能为空",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCtrlCategory,
	}
	errOrderSubmitTooFrequently = &hes.Error{
		Message:    "请勿重复提交订单",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCtrlCategory,
	}
)

func init() {
	ctrl := orderCtrl{}
	g := router.NewGroup("/orders")

	// 添加订单
	g.POST(
		"/v1",
		loadUserSession,
		shouldLogined,
		newTracker(cs.ActionOrderAdd),
		ctrl.add,
	)

	// 查看订单
	g.GET(
		"/v1",
		ctrl.list,
	)

	// 支付订单
	g.PATCH(
		"/v1/{sn}/pay",
		loadUserSession,
		shouldLogined,
		// 错误转换
		func(c *elton.Context) error {
			err := c.Next()
			if err == M.ErrSubmitTooFrequently {
				err = errOrderSubmitTooFrequently
			}
			return err
		},
		newTracker(cs.ActionOrderPay),
		middleware.NewConcurrentLimitWithDone([]string{
			"p:sn",
		}, time.Minute, ""),
		ctrl.pay,
	)

	g.GET(
		"/v1/statuses",
		ctrl.listStatus,
	)
	g.GET(
		"/v1/sub-order/statuses",
		ctrl.listSubOrderStatus,
	)
}

func (params listOrderParams) toConditions() (conditions []interface{}) {
	conds := queryConditions{}
	if params.Status != "" {
		conds.add("status = ?", params.Status)
	}
	if params.SN != "" {
		conds.add("sn = ?", params.SN)
	}

	return conds.toArray()
}

// add add order
func (orderCtrl) add(c *elton.Context) (err error) {
	params := addOrderParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	if len(params.Products) == 0 {
		err = errProductsIsEmpty
		return
	}
	us := getUserSession(c)
	subOrders := make([]service.SubOrder, len(params.Products))
	for index, prod := range params.Products {
		subOrders[index] = service.SubOrder{
			Product:      prod.ID,
			ProductCount: prod.Count,
		}
	}
	order, err := orderSrv.CreateWithSubOrders(us.GetID(), subOrders)
	if err != nil {
		return
	}
	c.Created(order)
	return
}

// list list order
func (orderCtrl) list(c *elton.Context) (err error) {
	params := listOrderParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	count := -1
	args := params.toConditions()
	queryParams := params.toPGQueryParams()
	if queryParams.Offset == 0 {
		count, err = orderSrv.Count(args...)
		if err != nil {
			return
		}
	}
	result, err := orderSrv.List(queryParams, args...)
	if err != nil {
		return
	}
	c.CacheMaxAge("1m")
	c.Body = &struct {
		Order []*service.Order `json:"order,omitempty"`
		Count int              `json:"count,omitempty"`
	}{
		result,
		count,
	}
	return
}

// listStatus list order status
func (orderCtrl) listStatus(c *elton.Context) (err error) {
	c.CacheMaxAge("5m")
	c.Body = &struct {
		Statuses service.OrderStatusInfoList `json:"statuses,omitempty"`
	}{
		orderSrv.ListOrderStatus(),
	}
	return
}

// listSubOrderStatus list sub order status
func (orderCtrl) listSubOrderStatus(c *elton.Context) (err error) {
	c.CacheMaxAge("5m")
	c.Body = &struct {
		Statuses service.SubOrderStatusInfoList `json:"statuses,omitempty"`
	}{
		orderSrv.ListSubOrderStatus(),
	}
	return
}

// pay pay order
func (orderCtrl) pay(c *elton.Context) (err error) {
	params := payOrderParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	sn := c.Param("sn")
	us := getUserSession(c)
	order, err := orderSrv.Pay(service.PayParams{
		UserID:    us.GetID(),
		PayAmount: params.PayAmount,
		SN:        sn,
		Source:    params.Source,
	})
	if err != nil {
		return
	}
	c.Body = order
	return
}

// toBeShipped set order to be shipped
func (orderCtrl) toBeShipped(c *elton.Context) (err error) {
	return
}
