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
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/vicanso/hes"
	"github.com/vicanso/origin/helper"
	"github.com/vicanso/origin/util"
	"gorm.io/gorm"
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
		PaySource string
	}
	// 创建订单参数
	CreateOrderParams struct {
		SubOrders []SubOrder
		// 订单总金额
		Amount              float64
		ReceiverName        string
		ReceiverMobile      string
		ReceiverBaseAddress string
		ReceiverAddress     string
	}
	// 订单状态时间线
	OrderStatusTimelineItem struct {
		CreatedAt  *time.Time  `json:"createdAt,omitempty"`
		Status     OrderStatus `json:"status,omitempty"`
		StatusDesc string      `json:"statusDesc,omitempty"`
	}
	OrderStatusTimeline []OrderStatusTimelineItem

	Orders []*Order
	// 订单记录
	Order struct {
		helper.Model

		Tx *gorm.DB `json:"-,omitempty" gorm:"-"`

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
		// 送货员
		Courier     uint   `json:"courier,omitempty" gorm:"index:idx_order_courier"`
		CourierName string `json:"courierName,omitempty" gorm:"-"`

		// 收货人
		ReceiverName   string `json:"receiverName,omitempty" gorm:"not null"`
		ReceiverMobile string `json:"receiverMobile,omitempty" gorm:"not null"`
		// 收货人地址（地址编码）
		ReceiverBaseAddress     string `json:"receiverBaseAddress,omitempty" gorm:"not null"`
		ReceiverBaseAddressDesc string `json:"receiverBaseAddressDesc,omitempty" grom:"-"`
		ReceiverAddress         string `json:"receiverAddress,omitempty" gorm:"not null"`

		// 时间
		PaidAt     *time.Time `json:"paidAt,omitempty"`
		DeliveryAt *time.Time `json:"deliveryAt,omitempty"`
		ReceivedAt *time.Time `json:"receivedAt,omitempty"`

		// TODO 添加source
		// 支付渠道
		PaySource string `json:"paySource,omitempty"`

		// 状态时间线
		StatusTimeline OrderStatusTimeline `json:"statusTimeline,omitempty"`
	}
	SubOrders []*SubOrder
	// 子订单记录
	SubOrder struct {
		helper.Model

		Tx *gorm.DB `json:"-" gorm:"-"`

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
	// 支付流水记录
	OrderPayment struct {
		helper.Model

		// 订单ID
		MainOrder uint `json:"mainOrder,omitempty" gorm:"unique_index:idx_order_payment_main_order"`
		// 用户ID（方便用户查询支付流水）
		UserID uint `json:"userID,omitempty" gorm:"index:idx_order_payment_user;not null"`
		// 支付渠道
		Source string `json:"source,omitempty" gorm:"not null"`
		// 支付金额
		PayAmount float64            `json:"payAmount,omitempty" gorm:"not null"`
		Status    OrderPaymentStatus `json:"status,omitempty"`
		Message   string             `json:"message,omitempty"`
	}
	// 订单派送记录
	OrderDelivery struct {
		helper.Model

		MainOrder uint   `json:"mainOrder,omitempty" gorm:"unique_index:idx_order_delivery_main_order"`
		UserID    uint   `json:"userID,omitempty" gorm:"index:idx_order_delivery_user;not null"`
		SN        string `json:"sn,omitempty"`
		Company   string `json:"company,omitempty"`
	}
	// OrderStatusSummary 订单状态概要
	OrderStatusSummary struct {
		Status     OrderStatus `json:"status,omitempty"`
		StatusDesc string      `json:"statusDesc,omitempty"`
		Count      int         `json:"count,omitempty"`
	}
	OrderSrv struct{}
)

const (
	errOrderCategory = "order"
)

const (
	OrderStatusUnknown OrderStatus = iota
	// 初始化
	OrderStatusInited
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
	SubOrderStatusUnknown SubOrderStatus = iota
	// 初始化
	SubOrderStatusInited
	// 待发货
	SubOrderStatusToBeShipped
	// 已发货
	SubOrderStatusShipped
	// 申请取消
	SubOrderStatusApplyCanceled
	// 取消
	SubOrderStatusCanceled
	// 申请退款
	SubOrderStatusApplyRefunds
	// 退款中
	SubOrderStatusRefunding
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
		SubOrderStatusShipped:       "已发货",
		SubOrderStatusApplyCanceled: "申请取消",
		SubOrderStatusCanceled:      "已取消",
		SubOrderStatusApplyRefunds:  "申请退款",
		SubOrderStatusRefunding:     "退款中",
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
		Message:    "订单金额异常，请重新刷新订单后再提交",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCategory,
	}
	errOrderProductInvalid = &hes.Error{
		Message:    "产品代码非法",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCategory,
	}
	errOrderOwnerInvalid = &hes.Error{
		Message:    "订单用户异常",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCategory,
	}
	errCanChangeToBeShipped = &hes.Error{
		Message:    "订单对应的子订单未处理完成，不可更新为待发货",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCategory,
	}
	errOrderCourierInvalid = &hes.Error{
		Message:    "该订单不属于你的配送订单或未分派配送员",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCategory,
	}
	errSubOrderNotMatch = &hes.Error{
		Message:    "子订单与订单不匹配",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCategory,
	}
	errCourierExists = &hes.Error{
		Message:    "该订单已分配派送员",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCategory,
	}
)

func init() {
	err := pgGetClient().AutoMigrate(
		&Order{},
		&SubOrder{},
		&OrderPayment{},
		&OrderDelivery{},
	)
	if err != nil {
		panic(err)
	}

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

func createOrderStatusTransferError(currentStatus, nextStatus OrderStatus) error {
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

func createSubOrderStatusTransferError(currentStatus, nextStatus SubOrderStatus) error {
	he := &hes.Error{
		Message:    fmt.Sprintf("子订单状态不能由%s至%s", currentStatus.String(), nextStatus.String()),
		Category:   errOrderCategory,
		StatusCode: http.StatusBadRequest,
	}
	return he
}

func containsSubOrderStatus(arr []SubOrderStatus, status SubOrderStatus) bool {
	found := false
	for _, item := range arr {
		if item == status {
			found = true
			break
		}
	}
	return found
}

func (timeline OrderStatusTimeline) Value() (driver.Value, error) {
	buf, err := json.Marshal(timeline)
	return string(buf), err
}

func (timeline *OrderStatusTimeline) Scan(input interface{}) error {
	switch value := input.(type) {
	case string:
		return json.Unmarshal([]byte(value), timeline)
	case []byte:
		return json.Unmarshal(value, timeline)
	default:
		return &hes.Error{
			Message:    "不支持的时间轴类型",
			Category:   errOrderCategory,
			StatusCode: http.StatusBadRequest,
		}
	}
}

// Add add status to timeline
func (timeline OrderStatusTimeline) Add(status OrderStatus) OrderStatusTimeline {
	now := time.Now()
	timeline = append(timeline, OrderStatusTimelineItem{
		CreatedAt:  &now,
		Status:     status,
		StatusDesc: status.String(),
	})
	return timeline
}

func (status OrderStatus) String() string {
	value, ok := orderStatusDict[status]
	if !ok {
		return ""
	}
	return value
}

func (status OrderStatus) Next() OrderStatus {
	nextStatus := OrderStatusUnknown
	switch status {
	case OrderStatusInited:
		nextStatus = OrderStatusPaymenting
	case OrderStatusPendingPayment:
		nextStatus = OrderStatusPaymenting
	case OrderStatusPaymenting:
		nextStatus = OrderStatusPaid
	case OrderStatusPaid:
		nextStatus = OrderStatusToBeShipped
	case OrderStatusToBeShipped:
		nextStatus = OrderStatusShipped
	case OrderStatusShipped:
		nextStatus = OrderStatusDone
	case OrderStatusDone:
		nextStatus = OrderStatusClosed
	}
	return nextStatus
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
	if status == nextStatus {
		err = &hes.Error{
			Message:    fmt.Sprintf("当前订单状态已是%s", status.String()),
			StatusCode: http.StatusBadRequest,
			Category:   errOrderCategory,
		}
		return
	}
	var allowStatuses []OrderStatus
	switch status {
	// 初始化成功的订单只能转向待支付
	case OrderStatusInited:
		allowStatuses = []OrderStatus{
			OrderStatusPendingPayment,
		}
		// 待支付 --> 支付中|已关闭
	case OrderStatusPendingPayment:
		allowStatuses = []OrderStatus{
			OrderStatusPaymenting,
			OrderStatusClosed,
		}
		// 支付中 --> 已支付|已关闭
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
		// 待发货 --> 已发货
	case OrderStatusToBeShipped:
		allowStatuses = []OrderStatus{
			OrderStatusShipped,
		}
		// 已发货 --> 已完成
	case OrderStatusShipped:
		allowStatuses = []OrderStatus{
			OrderStatusDone,
		}
		// 已完成 --> 已关闭
	case OrderStatusDone:
		allowStatuses = []OrderStatus{
			OrderStatusClosed,
		}
		// 已关闭订单不可更换状态
	case OrderStatusClosed:
		allowStatuses = []OrderStatus{}
	default:
		err = &hes.Error{
			Message:    fmt.Sprintf("异常状态[%s]", status.String()),
			Category:   errOrderCategory,
			StatusCode: http.StatusBadRequest,
		}
		return
	}

	// 如果下一状态非允许状态
	if !containsOrderStatus(allowStatuses, nextStatus) {
		err = createOrderStatusTransferError(status, nextStatus)
		return
	}
	return
}

// ValidateNext sub order validate next status
func (status SubOrderStatus) ValidateNext(nextStatus SubOrderStatus) (err error) {
	if status == nextStatus {
		err = &hes.Error{
			Message:    fmt.Sprintf("当前子订单状态已是%s", status.String()),
			StatusCode: http.StatusBadRequest,
			Category:   errOrderCategory,
		}
		return
	}
	var allowStatuses []SubOrderStatus
	switch status {
	// 初始化 --> 待发货|申请取消
	case SubOrderStatusInited:
		allowStatuses = []SubOrderStatus{
			SubOrderStatusToBeShipped,
			SubOrderStatusApplyCanceled,
		}
	// 待发货 --> 已发货
	case SubOrderStatusToBeShipped:
		allowStatuses = []SubOrderStatus{
			SubOrderStatusShipped,
		}
		// 已发货 --> 完成
	case SubOrderStatusShipped:
		allowStatuses = []SubOrderStatus{
			SubOrderStatusDone,
		}
		// 申请取消 --> 已取消
	case SubOrderStatusApplyCanceled:
		allowStatuses = []SubOrderStatus{
			SubOrderStatusCanceled,
		}
		// 已取消 --> 完成
	case SubOrderStatusCanceled:
		allowStatuses = []SubOrderStatus{
			SubOrderStatusDone,
		}
		// 申请退款 --> 退款中
	case SubOrderStatusApplyRefunds:
		allowStatuses = []SubOrderStatus{
			SubOrderStatusRefunding,
		}
		// 退款中 --> 完成
	case SubOrderStatusRefunding:
		allowStatuses = []SubOrderStatus{
			SubOrderStatusDone,
		}
		// 完成 --> 已关闭
	case SubOrderStatusDone:
		allowStatuses = []SubOrderStatus{
			SubOrderStatusClosed,
		}

	default:
		err = &hes.Error{
			Message:    fmt.Sprintf("异常状态[%d]", status),
			Category:   errOrderCategory,
			StatusCode: http.StatusBadRequest,
		}
		return
	}

	if !containsSubOrderStatus(allowStatuses, nextStatus) {
		err = createSubOrderStatusTransferError(status, nextStatus)
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

func (subOrder *SubOrder) BeforeCreate(_ *gorm.DB) error {
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

// UpdateStatus update sub order status
func (subOrder *SubOrder) UpdateStatus(status SubOrderStatus) (err error) {
	err = subOrder.Status.ValidateNext(status)
	if err != nil {
		return
	}
	db := subOrder.Tx
	if db == nil {
		db = pgGetClient()
	}

	// 保证当前的状态一致
	db = db.Model(subOrder).Where("status = ?", subOrder.Status).Updates(SubOrder{
		Status: status,
	})
	err = db.Error
	if err != nil {
		return
	}
	if db.RowsAffected != 1 {
		err = hes.New("更新子订单状态失败，该子订单当前状态已变化")
		return
	}
	subOrder.Status = status
	subOrder.StatusDesc = status.String()
	return
}

func (subOrders SubOrders) AfterFind(tx *gorm.DB) (err error) {
	for _, subOrder := range subOrders {
		err = subOrder.AfterFind(tx)
		if err != nil {
			return
		}
	}
	return
}

// TODO 针对子订单的状态
func (subOrder *SubOrder) AfterFind(_ *gorm.DB) (err error) {
	subOrder.StatusDesc = subOrder.Status.String()
	return
}

func (orders Orders) AfterFind(tx *gorm.DB) (err error) {
	for _, order := range orders {
		err = order.AfterFind(tx)
		if err != nil {
			return
		}
	}
	return
}

func (order *Order) BeforeCreate(_ *gorm.DB) (err error) {
	// 设置状态
	order.Status = OrderStatusInited
	timeline := make(OrderStatusTimeline, 0)
	timeline = timeline.Add(OrderStatusInited)
	order.StatusTimeline = timeline
	// order.StatusTimelineRaw = timeline.String()
	return
}

func (order *Order) FillAllStatusTimeline() {
	if len(order.StatusTimeline) != 0 {
		lastStatus := order.StatusTimeline[len(order.StatusTimeline)-1].Status
		// 最多只获取后5个状态
		for i := 0; i < 5; i++ {
			if lastStatus.Next() != OrderStatusUnknown {
				lastStatus = lastStatus.Next()
				order.StatusTimeline = append(order.StatusTimeline, OrderStatusTimelineItem{
					Status:     lastStatus,
					StatusDesc: lastStatus.String(),
				})
			}
		}
	}
}

func (order *Order) AfterFind(_ *gorm.DB) (err error) {
	order.StatusDesc = order.Status.String()
	order.CourierName, _ = userSrv.GetNameFromCache(order.Courier)

	// 收货地址不展示国家
	order.ReceiverBaseAddressDesc, _ = regionSrv.GetNameFromCache(order.ReceiverBaseAddress, 1)

	return
}

// UpdateStatus update order status
func (order *Order) UpdateStatus(status OrderStatus, updateDatas ...Order) (err error) {
	err = order.Status.ValidateNext(status)
	if err != nil {
		return
	}
	db := order.Tx
	if db == nil {
		db = pgGetClient()
	}

	timeline := order.StatusTimeline.Add(status)
	updateData := Order{}
	if len(updateDatas) != 0 {
		updateData = updateDatas[0]
	}
	updateData.Status = status
	updateData.StatusTimeline = timeline
	now := time.Now()
	switch status {
	case OrderStatusPaid:
		updateData.PaidAt = &now
	case OrderStatusShipped:
		updateData.DeliveryAt = &now
	case OrderStatusDone:
		updateData.ReceivedAt = &now
	}

	// 保证当前的状态一致
	db = db.Model(order).Where("status = ?", order.Status).Updates(updateData)
	err = db.Error
	if err != nil {
		return
	}
	if db.RowsAffected != 1 {
		err = hes.New("更新订单状态失败，该订单当前状态已变化")
		return
	}
	order.StatusTimeline = timeline
	order.Status = status
	order.StatusDesc = status.String()
	return
}

// ValidateCourier validate courier
func (order *Order) ValidateCourier(courier uint) error {
	if order.Courier != courier {
		return errOrderCourierInvalid
	}
	return nil
}

// ValidateOwner validate owner
func (order *Order) ValidateOwner(userID uint) error {
	if order.UserID != userID {
		return errOrderOwnerInvalid
	}
	return nil
}

func (payment *OrderPayment) BeforeCreate(_ *gorm.DB) (err error) {
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
func (srv *OrderSrv) CreateWithSubOrders(user uint, params CreateOrderParams) (order *Order, err error) {
	order = &Order{
		SN:                  srv.genSN(),
		UserID:              user,
		ReceiverName:        params.ReceiverName,
		ReceiverMobile:      params.ReceiverMobile,
		ReceiverBaseAddress: params.ReceiverBaseAddress,
		ReceiverAddress:     params.ReceiverAddress,
	}

	ids := make([]string, 0)
	for _, subOrder := range params.SubOrders {
		if subOrder.Product == 0 {
			err = errOrderProductInvalid
			return
		}
		id := strconv.Itoa(int(subOrder.Product))

		if !util.ContainsString(ids, id) {
			ids = append(ids, id)
		}
	}
	products, err := productSrv.List(PGQueryParams{
		Limit: len(ids),
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
		for _, subOrder := range params.SubOrders {
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
					err = tx.Create(&subOrder).Error
					if err != nil {
						return
					}
					amount += subOrder.ProductAmount
					payAmount += subOrder.ProductPayAmount
					// 已成功处理添加子订单，跳出循环
					break
				}
			}
			if !found {
				err = errOrderProductInvalid
				return
			}
		}
		// 如果应支付金额为0或者再客户端提交的金额不一致
		if payAmount == 0 || (payAmount != params.Amount) {
			err = errOrderAmountInValid
			return
		}

		err = tx.Model(order).Updates(Order{
			Status:    OrderStatusPendingPayment,
			Amount:    amount,
			PayAmount: payAmount,
		}).Error
		if err != nil {
			return
		}
		return
	})
	order.StatusDesc = order.Status.String()
	// 收货地址不展示国家
	order.ReceiverBaseAddressDesc, _ = regionSrv.GetNameFromCache(order.ReceiverBaseAddress, 1)
	return
}

// List list order
func (srv *OrderSrv) List(params PGQueryParams, args ...interface{}) (result Orders, err error) {
	result = make(Orders, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// COunt count order
func (srv *OrderSrv) Count(args ...interface{}) (count int64, err error) {
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

// FindSubOrdersByOrderID find sub orders by order id
func (srv *OrderSrv) FindSubOrdersByOrderID(orderID uint) (subOrders SubOrders, err error) {
	subOrders = make(SubOrders, 0)
	err = pgGetClient().Find(&subOrders, "main_order = ?", orderID).Error
	return
}

// FindSubOrdersByOrderIDList find sub order by order id list
func (srv *OrderSrv) FindSubOrdersByOrderIDList(orderIDList []uint) (subOrders SubOrders, err error) {
	subOrders = make(SubOrders, 0)
	err = pgGetClient().Find(&subOrders, "main_order IN (?)", orderIDList).Error
	return
}

// FindPaymentByOrderID find payment by order id
func (srv *OrderSrv) FindPaymentByOrderID(orderID uint) (orderPayment *OrderPayment, err error) {
	orderPayment = new(OrderPayment)
	err = pgGetClient().First(orderPayment, "main_order = ?", orderID).Error
	return
}

// FindSubOrderByID find sub order by id
func (srv *OrderSrv) FindSubOrderByID(subOrderID uint) (subOrder *SubOrder, err error) {
	subOrder = new(SubOrder)
	err = pgGetClient().First(subOrder, "id = ?", subOrderID).Error
	return
}

// Pay pay order
func (srv *OrderSrv) Pay(params PayParams) (order *Order, err error) {
	order, err = srv.FindBySN(params.SN)
	if err != nil {
		return
	}
	// TODO 如果账户对不上，有可能是攻击（正常账户不应该对不上)，可添加监控
	err = order.ValidateOwner(params.UserID)
	if err != nil {
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
				Source:    params.PaySource,
				PayAmount: params.PayAmount,
			}
			// TODO 添加支付流水
			err = tx.Create(orderPayment).Error
			if err != nil {
				return
			}
			// 同时更新父订单的支付渠道
			err = order.UpdateStatus(OrderStatusPaymenting, Order{
				PaySource: params.PaySource,
			})
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
	// 如果支付成功，则设置为成功
	// MOCK 设置为支付成功，失败则设置为失败

	paymentNextStatus := OrderPaymentStatusSuccess

	err = pgGetClient().Model(orderPayment).Updates(OrderPayment{
		Status: paymentNextStatus,
	}).Error
	// TODO 如果更新payment时失败，是否需要人手干预
	if err != nil {
		return
	}

	orderPayment.Status = paymentNextStatus
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

// ChangeCourier change order's courier
func (srv *OrderSrv) ChangeCourier(sn string, courier uint) (err error) {
	order, err := srv.FindBySN(sn)
	if err != nil {
		return
	}
	if order.Courier != 0 {
		err = errCourierExists
		return
	}
	err = srv.UpdateByID(order.ID, Order{
		Courier: courier,
	})
	if err != nil {
		return
	}
	return
}

// ToBeShipped order to be shipped
func (srv *OrderSrv) ToBeShipped(sn string, courier uint, subOrderID uint) (err error) {
	order, err := srv.FindBySN(sn)
	if err != nil {
		return
	}

	// 非送货员不可处理订单
	err = order.ValidateCourier(courier)
	if err != nil {
		return
	}

	err = order.Status.ValidateNext(OrderStatusToBeShipped)
	if err != nil {
		return
	}
	// 如果没有指定子订单，则表示整个订单
	if subOrderID == 0 {
		subOrders, e := srv.FindSubOrdersByOrderID(order.ID)
		if e != nil {
			err = e
			return
		}
		// 必须所有子订单都是待发货或已取消
		for _, item := range subOrders {
			if !containsSubOrderStatus([]SubOrderStatus{
				SubOrderStatusToBeShipped,
				SubOrderStatusCanceled,
			}, item.Status) {
				return errCanChangeToBeShipped
			}
		}
		err = order.UpdateStatus(OrderStatusToBeShipped)
		if err != nil {
			return
		}
		return
	}
	subOrder, err := srv.FindSubOrderByID(subOrderID)
	if err != nil {
		return
	}
	if subOrder.MainOrder != order.ID {
		err = errSubOrderNotMatch
		return
	}
	err = subOrder.UpdateStatus(SubOrderStatusToBeShipped)
	if err != nil {
		return
	}
	return
}

// Shipped set the order to shipped
func (srv *OrderSrv) Shipped(sn string, params OrderDelivery) (delivery *OrderDelivery, err error) {
	order, err := srv.FindBySN(sn)
	if err != nil {
		return
	}
	// 非送货员不可处理订单
	err = order.ValidateCourier(params.UserID)
	if err != nil {
		return
	}
	err = order.Status.ValidateNext(OrderStatusShipped)
	if err != nil {
		return
	}
	params.MainOrder = order.ID
	err = pgGetClient().Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Create(&params).Error
		if err != nil {
			return
		}
		order.Tx = tx
		err = tx.Model(&SubOrder{}).Where("main_order = ?", order.ID).Updates(SubOrder{
			Status: SubOrderStatusShipped,
		}).Error
		if err != nil {
			return
		}
		err = order.UpdateStatus(OrderStatusShipped)
		if err != nil {
			return
		}
		return
	})
	if err != nil {
		return
	}
	delivery = &params
	return
}

// changeStaus change order status
func (srv *OrderSrv) changeStaus(sn string, userID uint, nextStatus OrderStatus) (err error) {
	order, err := srv.FindBySN(sn)
	if err != nil {
		return
	}
	// 校验是否所有者
	err = order.ValidateOwner(userID)
	if err != nil {
		return
	}
	// 判断下一状态
	err = order.Status.ValidateNext(nextStatus)
	if err != nil {
		return
	}
	// 更新订单状态
	err = order.UpdateStatus(nextStatus)
	if err != nil {
		return
	}
	return
}

// Close close the order
func (srv *OrderSrv) Close(sn string, userID uint) (err error) {
	return srv.changeStaus(sn, userID, OrderStatusClosed)
}

// Finish finish the order
func (srv *OrderSrv) Finish(sn string, userID uint) (err error) {
	return srv.changeStaus(sn, userID, OrderStatusDone)
}

// ListStatusSummary list order status summary
func (srv *OrderSrv) ListStatusSummary(args ...interface{}) (summaryList []*OrderStatusSummary, err error) {
	db := pgGetClient().Model(&Order{})
	argsLen := len(args)
	if argsLen != 0 {
		if argsLen == 1 {
			db = db.Where(args[0])
		} else {
			db = db.Where(args[0], args[1:]...)
		}
	}
	orders := make(Orders, 0)
	err = db.Select("status").Find(&orders).Error
	if err != nil {
		return
	}
	summaryList = make([]*OrderStatusSummary, 0)
	for _, order := range orders {
		var found *OrderStatusSummary
		for _, summary := range summaryList {
			if summary.Status == order.Status {
				found = summary
				break
			}
		}
		if found != nil {
			found.Count++
		} else {
			summaryList = append(summaryList, &OrderStatusSummary{
				Status:     order.Status,
				StatusDesc: order.StatusDesc,
				Count:      1,
			})
		}
	}
	return
}
