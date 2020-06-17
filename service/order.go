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

package service

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/vicanso/hes"
	"github.com/vicanso/origin/helper"
	"github.com/vicanso/origin/util"
)

type (
	// 订单状态
	OrderStatus int
	// 子订单状态
	SubOrderStatus int
	// 支付订单状态
	OrderPaymentStatus int
	// 订单状态信息
	OrderStatusInfo struct {
		Name  string      `json:"name,omitempty"`
		Value OrderStatus `json:"value,omitempty"`
	}
	OrderStatusInfoList []*OrderStatusInfo
	SubOrderStatusInfo  struct {
		Name  string         `json:"name,omitempty"`
		Value SubOrderStatus `json:"value,omitempty"`
	}
	SubOrderStatusInfoList []*SubOrderStatusInfo
	// 支付参数
	PayParams struct {
		UserID    uint
		PayAmount float64
		SN        string
		Source    string
	}
	// 订单状态时间线
	OrderStatusTimelineItem struct {
		CreatedAt time.Time   `json:"createdAt,omitempty"`
		Status    OrderStatus `json:"status,omitempty"`
	}
	OrderStatusTimeline []OrderStatusTimelineItem

	Order struct {
		helper.Model

		// 编号
		SN string `json:"sn,omitempty" gorm:"not null;unique_index:idx_order_sn"`
		// 用户ID
		UserID uint `json:"userID,omitempty" gorm:"index:idx_order_user;not null"`
		// 总金额
		Amount float64 `json:"amount,omitempty" gorm:"not null"`
		// 支付金额
		PayAmount float64 `json:"payAmount,omitempty" gorm:"not null"`
		// 状态
		Status     OrderStatus `json:"status,omitempty" gorm:"index:idx_order_status"`
		StatusDesc string      `json:"statusDesc,omitempty" gorm:"-"`

		// 物流单号
		DeliverySN string `json:"deliverySN,omitempty" gorm:"index:idx_order_delivery_sn"`
		// 物流公司
		DeliveryCompany string `json:"deliveryCompany,omitempty"`

		// 收货人
		ReceiverName   string `json:"receiverName,omitempty"`
		ReceiverMobile string `json:"receiverMobile,omitempty"`
		// 收货人地址（地址编码）
		ReceiverBaseAddress string `json:"receiverBaseAddress,omitempty"`
		ReceiverAddress     string `json:"receiverAddress,omitempty"`

		// 时间
		PaidAt     *time.Time `json:"paidAt,omitempty"`
		DeliveryAt *time.Time `json:"deliveryAt,omitempty"`
		ReceivedAt *time.Time `json:"receivedAt,omitempty"`

		// 状态时间线
		StatusTimeline OrderStatusTimeline `json:"statusTimeline,omitempty" grom:"type:json"`
	}
	SubOrder struct {
		helper.Model

		MainOrder    uint    `json:"mainOrder,omitempty" gorm:"index:idx_sub_order_main_order"`
		Product      uint    `json:"product,omitempty" gorm:"not null"`
		ProductName  string  `json:"productName,omitempty" gorm:"not null"`
		ProductPrice float64 `json:"productPrice,omitempty" grom:"not null"`
		// 规格汇总
		ProductSpecsCount uint   `json:"productSpecsCount,omitempty" gorm:"not null"`
		ProductUnit       string `json:"productUnit,omitempty" gorm:"not null"`
		// 数量
		ProductCount uint `json:"productCount,omitempty" gorm:"not null"`
		// 金额
		ProductAmount float64 `json:"productAmount,omitempty" gorm:"not null"`
		// 支付金额
		ProductPayAmount float64 `json:"productPayAmount,omitempty" gorm:"not null"`
		// TODO 子订单状态
		// 状态
		Status SubOrderStatus `json:"status,omitempty" gorm:"index:idx_sub_order_status"`
		// 状态描述
		StatusDesc string `json:"statusDesc,omitempty" gorm:"-"`
	}
	OrderPayment struct {
		helper.Model

		// 订单ID
		MainOrder uint `json:"mainOrder,omitempty" gorm:"unique_index:idx_payment_order_main_order"`
		// 用户ID
		UserID uint `json:"userID,omitempty" gorm:"index:idx_payment_order_user;not null"`
		// 支付渠道
		Source string `json:"source,omitempty" gorm:"not null"`
		// 支付金额
		PayAmount float64            `json:"payAmount,omitempty" gorm:"not null"`
		Status    OrderPaymentStatus `json:"status,omitempty"`
		Message   string             `json:"message,omitempty"`
	}
	OrderSrv struct{}
)

const (
	errOrderCategory = "order"
)

const (
	// 初始化
	OrderStatusInited OrderStatus = iota + 1
	// 待支付
	OrderStatusPendingPayment
	// 正在支付
	OrderStatusPaymenting
	// 已支付
	OrderStatusPaid
	// 支付失败
	OrderStatusPayFail
	// 待发货
	OrderStatusToBeShipped
	// 已发货
	OrderStatusShipped
	// 已完成
	OrderStatusDone
	// 已关闭
	OrderStatusClosed
)

const (
	// 初始化
	SubOrderStatusInited SubOrderStatus = iota + 1
	// 待发货
	SubOrderStatusToBeShipped
	// 已发货
	SubOrderStatusShippe
	// 申请取消
	SubOrderStatusApplyCanceled
	// 取消
	SubOrderStatusCanceled
	// 退款中
	SubOrderStatusRefunds
	// 完成
	SubOrderStatusDone
	// 已关闭
	SubOrderStatusClosed
)

const (
	// 初始化
	OrderPaymentStatusInited OrderPaymentStatus = iota + 1
	// 失败
	OrderPaymentStatusFailure
	// 成功
	OrderPaymentStatusSuccess
)

var (
	orderStatusDict = map[OrderStatus]string{
		OrderStatusInited:         "初始化",
		OrderStatusPendingPayment: "待支付",
		OrderStatusPaymenting:     "支付中",
		OrderStatusPaid:           "已支付",
		OrderStatusPayFail:        "支付失败",
		OrderStatusToBeShipped:    "待发货",
		OrderStatusShipped:        "已发货",
		OrderStatusDone:           "已完成",
		OrderStatusClosed:         "已关闭",
	}
	orderStatusList    OrderStatusInfoList
	subOrderStatusDict = map[SubOrderStatus]string{
		SubOrderStatusInited:        "初始化",
		SubOrderStatusToBeShipped:   "待发货",
		SubOrderStatusShippe:        "已发货",
		SubOrderStatusApplyCanceled: "申请取消",
		SubOrderStatusCanceled:      "已取消",
		SubOrderStatusRefunds:       "退款中",
		SubOrderStatusDone:          "完成",
		SubOrderStatusClosed:        "已关闭",
	}
	subOrderStatusList SubOrderStatusInfoList
)

var (
	errOrderIdInvalid = &hes.Error{
		Message:    "订单ID不能为空",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCategory,
	}
	errOrderCountInvalid = &hes.Error{
		Message:    "购买数量必须大于等于1",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCategory,
	}
	errOrderAmountInValid = &hes.Error{
		Message:    "订单金额异常",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCategory,
	}
	errOrderProductInvalid = &hes.Error{
		Message:    "产品代码非法",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCategory,
	}
	errOrderIsInvalid = &hes.Error{
		Message:    "订单异常",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCategory,
	}
)

func init() {
	pgGetClient().AutoMigrate(
		&Order{},
		&SubOrder{},
		&OrderPayment{},
	)

	orderStatusList = make(OrderStatusInfoList, 0)
	for k, v := range orderStatusDict {
		orderStatusList = append(orderStatusList, &OrderStatusInfo{
			Name:  v,
			Value: k,
		})
	}
	sort.Slice(orderStatusList, func(i, j int) bool {
		return orderStatusList[i].Value < orderStatusList[j].Value
	})

	subOrderStatusList = make(SubOrderStatusInfoList, 0)
	for k, v := range subOrderStatusDict {
		subOrderStatusList = append(subOrderStatusList, &SubOrderStatusInfo{
			Name:  v,
			Value: k,
		})
	}
	sort.Slice(subOrderStatusList, func(i, j int) bool {
		return subOrderStatusList[i].Value < subOrderStatusList[j].Value
	})
}

func createStatusTransferError(currentStatus, nextStatus OrderStatus) error {
	he := &hes.Error{
		Message:    fmt.Sprintf("订单状态不能由%s至%s", currentStatus.String(), nextStatus.String()),
		Category:   errOrderCategory,
		StatusCode: http.StatusBadRequest,
	}
	return he
}

func containsOrderStatus(arr []OrderStatus, status OrderStatus) bool {
	found := false
	for _, item := range arr {
		if item == status {
			found = true
			break
		}
	}
	return found
}

func (status OrderStatus) String() string {
	value, ok := orderStatusDict[status]
	if !ok {
		return ""
	}
	return value
}

func (status SubOrderStatus) String() string {
	value, ok := subOrderStatusDict[status]
	if !ok {
		return ""
	}
	return value
}

// ValidateNext validate the status to next status
func (status OrderStatus) ValidateNext(nextStatus OrderStatus) (err error) {
	var allowStatuses []OrderStatus
	switch status {
	// 初始化成功的订单只能转向待支付
	case OrderStatusInited:
		allowStatuses = []OrderStatus{
			OrderStatusPendingPayment,
		}
		// 待支付 --> 正在支付|已关闭
	case OrderStatusPendingPayment:
		allowStatuses = []OrderStatus{
			OrderStatusPaymenting,
			OrderStatusClosed,
		}
		// 正在支付 --> 已支付|已关闭
	case OrderStatusPaymenting:
		allowStatuses = []OrderStatus{
			OrderStatusPaid,
			OrderStatusClosed,
		}
		// 已支付 --> 待发货
	case OrderStatusPaid:
		allowStatuses = []OrderStatus{
			OrderStatusToBeShipped,
		}
		// 支付失败 --> 已关闭
	case OrderStatusPayFail:
		allowStatuses = []OrderStatus{
			OrderStatusClosed,
		}
		// 已发货 --> 已完成
	case OrderStatusToBeShipped:
		allowStatuses = []OrderStatus{
			OrderStatusDone,
		}
		// 已完成 --> 已关闭
	case OrderStatusDone:
		allowStatuses = []OrderStatus{
			OrderStatusClosed,
		}
	default:
		err = &hes.Error{
			Message:    fmt.Sprintf("异常状态[%d]", status),
			Category:   errOrderCategory,
			StatusCode: http.StatusBadRequest,
		}
		return
	}

	// 如果下一状态非允许状态
	if !containsOrderStatus(allowStatuses, nextStatus) {
		err = createStatusTransferError(status, nextStatus)
		return
	}
	return
}

// CheckValid check sub order is valid
func (subOrder *SubOrder) CheckValid() error {
	if subOrder.ProductCount == 0 {
		return errOrderCountInvalid
	}
	if subOrder.MainOrder == 0 {
		return errOrderIdInvalid
	}
	return nil
}

func (subOrder *SubOrder) BeforeCreate() error {
	err := subOrder.CheckValid()
	if err != nil {
		return err
	}
	subOrder.Status = SubOrderStatusInited
	subOrder.ProductAmount = subOrder.ProductPrice * float64(subOrder.ProductCount)
	// 支付价格暂时无优惠
	subOrder.ProductPayAmount = subOrder.ProductAmount
	return nil
}

// TODO 针对子订单的状态
func (subOrder *SubOrder) AfterFind() (err error) {
	subOrder.StatusDesc = subOrder.Status.String()
	return
}

func (order *Order) BeforeCreate() (err error) {
	// 设置状态
	order.Status = OrderStatusInited
	timeline := make(OrderStatusTimeline, 0)
	timeline = append(timeline, OrderStatusTimelineItem{
		Status:    OrderStatusInited,
		CreatedAt: time.Now(),
	})
	order.StatusTimeline = timeline
	return
}

func (order *Order) AfterFind() (err error) {
	order.StatusDesc = order.Status.String()
	return
}

// UpdateStatus update order status
func (order *Order) UpdateStatus(status OrderStatus) (err error) {
	err = order.Status.ValidateNext(status)
	if err != nil {
		return
	}
	// 保证当前的状态一致
	db := pgGetClient().Model(order).Where("status = ?", order.Status).Update(Order{
		Status: status,
	})
	err = db.Error
	if err != nil {
		return
	}
	if db.RowsAffected != 1 {
		err = hes.New("更新订单状态失败，该订单当前状态已变化")
		return
	}
	fmt.Println(db.RowsAffected)
	order.Status = status
	order.StatusDesc = status.String()
	return
}

func (payment *OrderPayment) BeforeCreate() (err error) {
	// 设置初始化状态
	payment.Status = OrderPaymentStatusInited
	return
}

func (srv *OrderSrv) genSN() string {
	return util.GenUlid()
}

func (srv *OrderSrv) createByID(id uint) *Order {
	order := &Order{}
	order.Model.ID = id
	return order
}

// ListOrderStatus list the status of order
func (srv *OrderSrv) ListOrderStatus() OrderStatusInfoList {
	return orderStatusList
}

// ListSubOrderStatus list the status of sub order
func (srv *OrderSrv) ListSubOrderStatus() SubOrderStatusInfoList {
	return subOrderStatusList
}

// CreateWithSubOrders create order with sub orders
func (srv *OrderSrv) CreateWithSubOrders(user uint, data []SubOrder) (order *Order, err error) {
	order = &Order{
		SN:     srv.genSN(),
		UserID: user,
	}
	subOrders := make([]*SubOrder, len(data))
	for i, order := range data {
		subOrders[i] = &order
	}
	ids := make([]string, 0)
	for _, subOrder := range subOrders {
		id := strconv.Itoa(int(subOrder.Product))
		if !util.ContainsString(ids, id) {
			ids = append(ids, id)
		}
	}
	products, err := productSrv.List(PGQueryParams{
		Limit:  len(ids),
		Fields: "-catalog,categories",
	}, "id IN (?)", ids)
	if err != nil {
		return
	}

	for _, p := range products {
		err = p.CheckAvailable()
		if err != nil {
			return
		}
	}

	err = pgGetClient().Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Create(order).Error
		if err != nil {
			return
		}
		var amount, payAmount float64
		for _, subOrder := range subOrders {
			found := false
			for _, p := range products {
				// 订单记录产品当前信息，避免产品更新后，信息不符合
				if subOrder.Product == p.ID {
					found = true
					subOrder.ProductName = p.Name
					subOrder.ProductPrice = p.Price
					subOrder.ProductSpecsCount = p.Specs * subOrder.ProductCount
					subOrder.ProductUnit = p.Unit
					subOrder.MainOrder = order.ID
					err = tx.Create(subOrder).Error
					if err != nil {
						return
					}
					amount += subOrder.ProductAmount
					payAmount += subOrder.ProductPayAmount
				}
			}
			if !found {
				err = errOrderProductInvalid
				return
			}
		}
		if amount == 0 || payAmount == 0 {
			err = errOrderAmountInValid
			return
		}

		err = tx.Model(order).Update(&Order{
			Status:    OrderStatusPendingPayment,
			Amount:    amount,
			PayAmount: payAmount,
		}).Error
		if err != nil {
			return
		}
		return
	})
	return
}

// List list order
func (srv *OrderSrv) List(params PGQueryParams, args ...interface{}) (result []*Order, err error) {
	result = make([]*Order, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// COunt count order
func (srv *OrderSrv) Count(args ...interface{}) (count int, err error) {
	return pgCount(&Order{}, args...)
}

// FindBySN find order by sn
func (srv *OrderSrv) FindBySN(sn string) (order *Order, err error) {
	order = new(Order)
	err = pgGetClient().First(order, "sn = ?", sn).Error
	return
}

// UpdateByID update order by id
func (srv *OrderSrv) UpdateByID(id uint, order Order) (err error) {
	err = pgGetClient().Model(srv.createByID(id)).Updates(order).Error
	return
}

// FindPaymentByOrderID find payment by order id
func (srv *OrderSrv) FindPaymentByOrderID(orderID uint) (orderPayment *OrderPayment, err error) {
	orderPayment = new(OrderPayment)
	err = pgGetClient().First(orderPayment, "main_order = ?", orderID).Error
	return
}

// Pay pay order
func (srv *OrderSrv) Pay(params PayParams) (order *Order, err error) {
	order, err = srv.FindBySN(params.SN)
	if err != nil {
		return
	}
	// TODO 如果账户对不上，有可能是攻击（正常账户不应该对不上)，可添加监控
	if params.UserID != order.UserID {
		err = errOrderIsInvalid
		return
	}
	if params.PayAmount != order.PayAmount {
		err = &hes.Error{
			Message:    fmt.Sprintf("支付金额错误，应支付:%.2f", order.Amount),
			StatusCode: http.StatusBadRequest,
			Category:   errOrderCategory,
		}
		return
	}
	// TODO 当支付时，是否预生成支付流水（支付渠道等）
	var orderPayment *OrderPayment
	// 如果是待支付，增加记录当前准备支付
	if order.Status == OrderStatusPendingPayment {
		err = pgGetClient().Transaction(func(tx *gorm.DB) (err error) {
			orderPayment = &OrderPayment{
				MainOrder: order.ID,
				UserID:    order.UserID,
				Source:    params.Source,
				PayAmount: params.PayAmount,
			}
			// TODO 添加支付流水
			err = tx.Create(orderPayment).Error
			if err != nil {
				return
			}
			err = order.UpdateStatus(OrderStatusPaymenting)
			if err != nil {
				return
			}
			return
		})
		if err != nil {
			return
		}
	} else {
		orderPayment, err = srv.FindPaymentByOrderID(order.ID)
		if err != nil {
			return
		}
	}

	// 判断当前订单是否可流转为已支付
	err = order.Status.ValidateNext(OrderStatusPaid)
	if err != nil {
		return
	}

	// TODO 根据orderPayment去支付
	// MOCK 设置为支付成功
	orderPayment.Status = OrderPaymentStatusSuccess
	nextStatus := OrderStatusPayFail
	if orderPayment.Status == OrderPaymentStatusSuccess {
		nextStatus = OrderStatusPaid
	}

	// 成功支付则设置订单为已支付
	err = order.UpdateStatus(nextStatus)
	if err != nil {
		return
	}

	return
}

func (srv *OrderSrv) ToBeShipped(sn string) (order *Order, err error) {
	order, err = srv.FindBySN(sn)
	if err != nil {
		return
	}
	return
}
