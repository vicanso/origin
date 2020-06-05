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
	"net/http"
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
	Order          struct {
		helper.Model

		// 编号
		SN string `json:"sn,omitempty" gorm:"not null;unique_index:idx_order_sn"`
		// 用户ID
		UserID uint `json:"userID,omitempty" gorm:"index:idx_order_user"`
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
	}
	SubOrder struct {
		helper.Model

		Order        uint    `json:"order,omitempty" gorm:"index:idx_sub_order_order"`
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
	OrderSrv struct{}
)

const (
	errOrderCatgory = "order"
)

const (
	// 初始化
	OrderStatusInited OrderStatus = iota + 1
	// 待支付
	OrderStatusPendingPayment
	// 已支付
	OrderStatusPaid
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

var orderStatusDict = map[OrderStatus]string{
	OrderStatusInited:         "初始化",
	OrderStatusPendingPayment: "待支付",
	OrderStatusPaid:           "已支付",
	OrderStatusToBeShipped:    "待发货",
	OrderStatusShipped:        "已发货",
	OrderStatusDone:           "已完成",
	OrderStatusClosed:         "已关闭",
}
var subOrderStatusDict = map[SubOrderStatus]string{
	SubOrderStatusInited:        "初始化",
	SubOrderStatusApplyCanceled: "申请取消",
	SubOrderStatusCanceled:      "已取消",
	SubOrderStatusRefunds:       "退款中",
	SubOrderStatusDone:          "完成",
	SubOrderStatusClosed:        "已关闭",
}

var (
	errOrderIdInvalid = &hes.Error{
		Message:    "订单ID不能为空",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCatgory,
	}
	errOrderCountInvalid = &hes.Error{
		Message:    "购买数量必须大于等于1",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCatgory,
	}
	errOrderAmountInValid = &hes.Error{
		Message:    "订单金额异常",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCatgory,
	}
	errOrderProductInvalid = &hes.Error{
		Message:    "产品代码非法",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCatgory,
	}
)

func init() {
	pgGetClient().AutoMigrate(&Order{}).
		AutoMigrate(&SubOrder{})
}

// CheckValid check sub order is valid
func (subOrder *SubOrder) CheckValid() error {
	if subOrder.ProductCount <= 0 {
		return errOrderCountInvalid
	}
	if subOrder.Order == 0 {
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
	value, ok := subOrderStatusDict[subOrder.Status]
	if ok {
		subOrder.StatusDesc = value
	}
	return
}

func (order *Order) BeforeCreate() (err error) {
	// 设置状态
	order.Status = OrderStatusInited
	return
}

func (order *Order) AfterFind() (err error) {
	value, ok := orderStatusDict[order.Status]
	if ok {
		order.StatusDesc = value
	}
	return
}

func (srv *OrderSrv) genSN() string {
	return util.GenUlid()
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
	products, err := productSrv.List(helper.PGQueryParams{
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
					subOrder.Order = order.ID
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
