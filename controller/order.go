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
	"strconv"
	"strings"
	"time"

	"github.com/vicanso/elton"
	M "github.com/vicanso/elton/middleware"
	"github.com/vicanso/hes"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/middleware"
	"github.com/vicanso/origin/router"
	"github.com/vicanso/origin/service"
	"github.com/vicanso/origin/util"
	"github.com/vicanso/origin/validate"
	"gorm.io/gorm"
)

type (
	orderCtrl struct{}

	addOrderParams struct {
		Products []struct {
			ProductID uint    `json:"productID,omitempty" validate:"xOrderProductID"`
			Count     uint    `json:"count,omitempty" validate:"xOrderProductCount"`
			Price     float64 `json:"price,omitempty" validate:"xProductPrice"`
		} `json:"products,omitempty"`
		Amount              float64 `json:"amount,omitempty" validate:"required"`
		ReceiverName        string  `json:"receiverName,omitempty"`
		ReceiverMobile      string  `json:"receiverMobile,omitempty" validate:"xMobile"`
		ReceiverBaseAddress string  `json:"receiverBaseAddress,omitempty" validate:"xBaseAddress"`
		ReceiverAddress     string  `json:"receiverAddress,omitempty" validate:"xAddress"`
	}
	// 支付参数
	payOrderParams struct {
		PayAmount float64 `json:"payAmount,omitempty" validate:"xOrderAmount"`
		PaySource string  `json:"paySource,omitempty" validate:"xPaySource"`
	}
	// 订单待发货参数
	toBeShippedOrderParams struct {
		SubOrder uint `json:"subOrder,omitempty"`
	}
	// 订单发货参数
	shippedOrderParams struct {
		DeliverySN      string `json:"deliverySN,omitempty" validate:"xOrderDeliverySN"`
		DeliveryCompany string `json:"deliveryCompany,omitempty" validate:"xOrderDeliveryCompnay"`
	}
	// 修改送货人参数
	changeOrderDelivererParams struct {
		Deliverer uint `json:"deliverer,omitempty" validate:"xOrderDeliverer"`
	}

	listOrderParams struct {
		listParams

		Status            string    `json:"status,omitempty" validate:"omitempty,xOrderStatus"`
		Statuses          string    `json:"statuses,omitempty"`
		SN                string    `json:"sn,omitempty" validate:"omitempty,xOrderSN"`
		Begin             time.Time `json:"begin,omitempty"`
		End               time.Time `json:"end,omitempty"`
		User              string    `json:"user,omitempty" validate:"omitempty,xOrderUser"`
		Deliverer         string    `json:"deliverer,omitempty" validate:"omitempty,xOrderDeliverer"`
		IncludingSubOrder string    `json:"includingSubOrder,omitempty"`
		// 未派送订单
		NoDelivery string `json:"noDelivery,omitempty"`
	}
	// listOrderResp 订单列表响应
	listOrderResp struct {
		Orders    service.Orders    `json:"orders,omitempty"`
		SubOrders service.SubOrders `json:"subOrders,omitempty"`
		Count     int64             `json:"count,omitempty"`
	}

	// updateLocationParams 更新定位参数(暂时不可能出现0, 0的定位，因此设置为required)
	updateLocationParams struct {
		Latitude  float64 `json:"latitude,omitempty" validate:"xLatitude,required"`
		Longitude float64 `json:"longitude,omitempty" validate:"xLongitude,required"`
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
	orderUpdateLimit := elton.Compose(
		// 错误转换
		func(c *elton.Context) error {
			err := c.Next()
			if err == M.ErrSubmitTooFrequently {
				err = errOrderSubmitTooFrequently
			}
			return err
		},
		middleware.NewConcurrentLimitWithDone([]string{
			"p:sn",
		}, time.Minute, ""),
	)

	// 添加订单
	g.POST(
		"/v1",
		loadUserSession,
		shouldBeLogined,
		newTracker(cs.ActionOrderAdd),
		// 根据客户端提交的token限制同时提交
		middleware.NewConcurrentLimitWithDone([]string{
			"token",
		}, time.Minute, ""),
		ctrl.add,
	)

	// 查看订单
	g.GET(
		"/v1",
		loadUserSession,
		shouldBeLogined,
		checkMarketingGroup,
		ctrl.list,
	)
	// 查看我的订单
	g.GET(
		"/v1/mine",
		loadUserSession,
		shouldBeLogined,
		ctrl.listMine,
	)
	// 我的订单概要
	g.GET(
		"/v1/mine/summary",
		loadUserSession,
		shouldBeLogined,
		ctrl.listMineSummary,
	)
	// 查询派送订单
	g.GET(
		"/v1/my-deliveries",
		loadUserSession,
		shouldBeLogined,
		checkLogisticsGroup,
		ctrl.listDeliveryOrder,
	)
	// 查询未分派订单
	g.GET(
		"/v1/no-delivery",
		loadUserSession,
		shouldBeLogined,
		checkLogisticsGroup,
		ctrl.listNoDelivery,
	)

	// TODO 查订单详情是否只允许本人或管理人员查
	g.GET(
		"/v1/{sn}",
		loadUserSession,
		shouldBeLogined,
		ctrl.detail,
	)

	// 支付订单
	g.PATCH(
		"/v1/{sn}/pay",
		loadUserSession,
		shouldBeLogined,
		newTracker(cs.ActionOrderPay),
		orderUpdateLimit,
		ctrl.pay,
	)
	// 关闭订单
	g.PATCH(
		"/v1/{sn}/close",
		loadUserSession,
		shouldBeLogined,
		newTracker(cs.ActionOrderClose),
		orderUpdateLimit,
		ctrl.close,
	)
	// 结束订单
	g.PATCH(
		"/v1/{sn}/finish",
		loadUserSession,
		shouldBeLogined,
		newTracker(cs.ActionOrderFinish),
		orderUpdateLimit,
		ctrl.finish,
	)

	// 分派送货员
	g.PATCH(
		"/v1/{sn}/assign-deliverer",
		loadUserSession,
		shouldBeLogined,
		newTracker(cs.ActionOrderChangeDeliverer),
		checkMarketingGroup,
		orderUpdateLimit,
		ctrl.changeDeliverer,
	)
	// 抢接派单
	g.PATCH(
		"/v1/{sn}/assign-deliverer-to-me",
		loadUserSession,
		shouldBeLogined,
		newTracker(cs.ActionOrderChangeDelivererToMe),
		checkLogisticsGroup,
		orderUpdateLimit,
		ctrl.changeDelivererToMe,
	)
	// 订单设置为待发货
	g.PATCH(
		"/v1/{sn}/to-be-shipped",
		loadUserSession,
		shouldBeLogined,
		newTracker(cs.ActionOrderToBeShipped),
		checkLogisticsGroup,
		orderUpdateLimit,
		ctrl.toBeShipped,
	)
	// 订单设置为已发货
	g.PATCH(
		"/v1/{sn}/ship",
		loadUserSession,
		shouldBeLogined,
		newTracker(cs.ActionOrderShipped),
		checkLogisticsGroup,
		orderUpdateLimit,
		ctrl.shipped,
	)
	// 更新正在派送订单的location
	g.PATCH(
		"/v1/delivering/loaction",
		loadUserSession,
		shouldBeLogined,
		newTracker(cs.ActionOrderUpdateDeliveryLocation),
		checkLogisticsGroup,
		ctrl.updateDeliveringLocation,
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
	// 未分配派送订单
	if params.NoDelivery != "" {
		params.Deliverer = "0"
		params.Status = strconv.Itoa(int(service.OrderStatusPaid))
	}

	if params.Status != "" {
		conds.add("status = ?", params.Status)
	}
	if params.SN != "" {
		conds.add("sn = ?", params.SN)
	}

	if !params.Begin.IsZero() {
		conds.add("created_at >= ?", util.FormatTime(params.Begin))
	}
	if !params.End.IsZero() {
		conds.add("created_at <= ?", util.FormatTime(params.End))
	}
	if params.User != "" {
		id, _ := strconv.Atoi(params.User)
		conds.add("user_id = ?", id)
	}

	if params.Deliverer != "" {
		id, _ := strconv.Atoi(params.Deliverer)
		conds.add("deliverer = ?", id)
	}
	if params.Statuses != "" {
		conds.add("status in (?)", strings.Split(params.Statuses, ","))
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
			Product:      prod.ProductID,
			ProductCount: prod.Count,
			ProductPrice: prod.Price,
		}
	}
	order, err := orderSrv.CreateWithSubOrders(us.GetID(), service.CreateOrderParams{
		Amount:              params.Amount,
		SubOrders:           subOrders,
		ReceiverName:        params.ReceiverName,
		ReceiverMobile:      params.ReceiverMobile,
		ReceiverBaseAddress: params.ReceiverBaseAddress,
		ReceiverAddress:     params.ReceiverAddress,
	})
	if err != nil {
		return
	}
	c.Created(order)
	return
}

func (orderCtrl) listOrder(params listOrderParams) (resp *listOrderResp, err error) {
	count := int64(-1)
	args := params.toConditions()
	queryParams := params.toPGQueryParams()
	if queryParams.Offset == 0 {
		count, err = orderSrv.Count(args...)
		if err != nil {
			return
		}
	}
	orders, err := orderSrv.List(queryParams, args...)
	if err != nil {
		return
	}
	var subOrders service.SubOrders
	if params.IncludingSubOrder == "true" && len(orders) != 0 {
		idList := make([]uint, len(orders))
		for index, order := range orders {
			idList[index] = order.ID
		}
		subOrders, err = orderSrv.FindSubOrdersByOrderIDList(idList)
		if err != nil {
			return
		}
	}
	resp = &listOrderResp{
		Orders:    orders,
		SubOrders: subOrders,
		Count:     count,
	}
	return
}

// list list order
func (ctrl orderCtrl) list(c *elton.Context) (err error) {
	params := listOrderParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	resp, err := ctrl.listOrder(params)
	if err != nil {
		return
	}
	c.Body = resp
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

// detail get the order detail
func (orderCtrl) detail(c *elton.Context) (err error) {
	sn := c.Param("sn")
	order, err := orderSrv.FindBySN(sn)
	if err != nil {
		return
	}
	order.FillAllStatusTimeline()
	subOrders, err := orderSrv.FindSubOrdersByOrderID(order.ID)
	if err != nil {
		return
	}
	var payment *service.OrderPayment
	// 初始化订单或待支付的则不需要查询payment记录
	if order.Status != service.OrderStatusInited &&
		order.Status != service.OrderStatusPendingPayment {
		payment, err = orderSrv.FindPaymentByOrderID(order.ID)
		if err == gorm.ErrRecordNotFound {
			payment = nil
			err = nil
		}
		if err != nil {
			return
		}
	}
	var delivery *service.OrderDelivery
	if order.Status == service.OrderStatusShipped {
		delivery, err = orderSrv.FindDeliveryByOrderID(order.ID)
		if err == gorm.ErrRecordNotFound {
			delivery = nil
			err = nil
		}
		if err != nil {
			return
		}
	}

	if err != nil {
		return
	}
	c.Body = &struct {
		Order     *service.Order         `json:"order,omitempty"`
		SubOrders service.SubOrders      `json:"subOrders,omitempty"`
		Payment   *service.OrderPayment  `json:"payment,omitempty"`
		Delivery  *service.OrderDelivery `json:"delivery,omitempty"`
	}{
		order,
		subOrders,
		payment,
		delivery,
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
		PaySource: params.PaySource,
	})
	if err != nil {
		return
	}
	c.Body = order
	return
}

// changeDeliverer change deliverer
func (orderCtrl) changeDeliverer(c *elton.Context) (err error) {
	params := changeOrderDelivererParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	// 如果由后台调整，则可强制调整派送员
	err = orderSrv.ChangeDeliverer(c.Param("sn"), params.Deliverer, true)
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// changeDelivererToMe change deliverer to me
func (orderCtrl) changeDelivererToMe(c *elton.Context) (err error) {
	us := getUserSession(c)
	err = orderSrv.ChangeDeliverer(c.Param("sn"), us.GetID(), false)
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// toBeShipped set order to be shipped
func (orderCtrl) toBeShipped(c *elton.Context) (err error) {
	params := toBeShippedOrderParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	us := getUserSession(c)
	err = orderSrv.ToBeShipped(c.Param("sn"), us.GetID(), params.SubOrder)
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// shipped set order to shipped
func (orderCtrl) shipped(c *elton.Context) (err error) {
	params := shippedOrderParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}

	us := getUserSession(c)
	delivery, err := orderSrv.Shipped(c.Param("sn"), service.OrderDelivery{
		UserID:  us.GetID(),
		SN:      params.DeliverySN,
		Company: params.DeliveryCompany,
	})
	if err != nil {
		return
	}
	c.Body = delivery
	return
}

// close set the order to closed
func (orderCtrl) close(c *elton.Context) (err error) {
	us := getUserSession(c)
	err = orderSrv.Close(c.Param("sn"), us.GetID())
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// finish set the order to done
func (orderCtrl) finish(c *elton.Context) (err error) {
	us := getUserSession(c)
	err = orderSrv.Finish(c.Param("sn"), us.GetID())
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// listDeliveryOrder list the delivery order
func (ctrl orderCtrl) listDeliveryOrder(c *elton.Context) (err error) {
	params := listOrderParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	us := getUserSession(c)
	// 避免通过user参数查询
	params.User = ""
	params.Deliverer = strconv.Itoa(int(us.GetID()))
	resp, err := ctrl.listOrder(params)
	if err != nil {
		return
	}
	c.Body = resp
	return
}

// listMine list my orders
func (ctrl orderCtrl) listMine(c *elton.Context) (err error) {
	params := listOrderParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	us := getUserSession(c)
	params.User = strconv.FormatInt(int64(us.GetID()), 10)
	resp, err := ctrl.listOrder(params)
	if err != nil {
		return
	}
	c.Body = resp
	return
}

// listMineSummary
func (ctrl orderCtrl) listMineSummary(c *elton.Context) (err error) {
	params := listOrderParams{}
	us := getUserSession(c)
	params.User = strconv.FormatInt(int64(us.GetID()), 10)
	// 如果未指定则查最近一个月的订单状态
	if params.Begin.IsZero() {
		params.Begin = time.Now().AddDate(0, -1, 0)
	}
	summaryList, err := orderSrv.ListStatusSummary(params.toConditions()...)
	if err != nil {
		return
	}
	c.Body = &struct {
		OrderStatusSummaries []*service.OrderStatusSummary `json:"orderStatusSummaries,omitempty"`
	}{
		summaryList,
	}
	return
}

// listNoDelivery
func (ctrl orderCtrl) listNoDelivery(c *elton.Context) (err error) {
	params := listOrderParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	params.NoDelivery = "1"
	resp, err := ctrl.listOrder(params)
	if err != nil {
		return
	}
	c.Body = resp
	return
}

// updateDeliveringLocation 更新正在派送的定位信息
func (orderCtrl) updateDeliveringLocation(c *elton.Context) (err error) {
	params := updateLocationParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	us := getUserSession(c)
	err = orderSrv.UpdateDeliveringLocation(us.GetID(), service.LocationTimelineItem{
		Latitude:  params.Latitude,
		Longitude: params.Longitude,
	})
	if err != nil {
		return
	}
	c.NoContent()
	return
}
