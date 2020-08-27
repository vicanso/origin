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
	"encoding/json"
	"sync"
	"time"

	"github.com/vicanso/origin/helper"
	"github.com/vicanso/origin/util"
)

type (
	// OrderCommission 订单佣金
	OrderCommission struct {
		helper.Model

		UserID  uint   `json:"userID,omitempty" gorm:"index:idx_order_commission_user_id;not null"`
		OrderSN string `json:"orderSN,omitempty" gorm:"unique_index:idx_order_commission_order_sn;not null"`
		// 推荐人
		Recommender      uint    `json:"recommender,omitempty" gorm:"index:idx_order_commission_recommender;not null"`
		PayAmount        float64 `json:"payAmount,omitempty" gorm:"not null"`
		CommissionAmount float64 `json:"commissionAmount,omitempty" gorm:"not null"`
	}

	OrderCommissions []*OrderCommission

	// OrderCommissionSrv 订单佣金
	OrderCommissionSrv struct{}

	OrderCommissionConfig struct {
		Group string  `json:"group,omitempty"`
		Ratio float64 `json:"ratio,omitempty"`
	}
	// OrderCommissionConfigs 订单佣金配置
	OrderCommissionConfigs struct {
		sync.RWMutex
		configs []*OrderCommissionConfig
	}
)

var (
	defaultOrderCommissionConfigs = new(OrderCommissionConfigs)
)

const (
	orderCommissionAllGroup = "*"
)

func init() {
	err := pgGetClient().AutoMigrate(
		&OrderCommission{},
	)
	if err != nil {
		panic(err)
	}
}

func (orderCommissionConfigs *OrderCommissionConfigs) Set(items []string) {
	confs := make([]*OrderCommissionConfig, 0)
	for _, item := range items {
		conf := &OrderCommissionConfig{}
		// 忽略出错
		_ = json.Unmarshal([]byte(item), conf)
		if conf.Group != "" && conf.Ratio >= 0 {
			confs = append(confs, conf)
		}
	}
	orderCommissionConfigs.Lock()
	defer orderCommissionConfigs.Unlock()
	orderCommissionConfigs.configs = confs
}

func (orderCommissionConfigs *OrderCommissionConfigs) Get(group string) (ratio float64) {
	orderCommissionConfigs.RLock()
	defer orderCommissionConfigs.RUnlock()
	for _, item := range orderCommissionConfigs.configs {
		if item.Group == group {
			ratio = item.Ratio
		}
	}
	return
}

func (srv *OrderCommissionSrv) createOrUpdate(order *Order) (err error) {
	if order.Recommender == 0 {
		return
	}
	ratio := defaultOrderCommissionConfigs.Get(orderCommissionAllGroup)
	if ratio == 0 {
		return
	}
	orderCommission := &OrderCommission{
		UserID:           order.UserID,
		OrderSN:          order.SN,
		Recommender:      order.Recommender,
		PayAmount:        order.PayAmount,
		CommissionAmount: order.PayAmount * ratio,
	}
	err = pgGetClient().FirstOrCreate(orderCommission, OrderCommission{
		OrderSN: order.SN,
	}).Error
	return
}

func (srv *OrderCommissionSrv) FindBySN(sn string) (orderCommission *OrderCommission, err error) {
	orderCommission = new(OrderCommission)
	err = pgGetClient().First(orderCommission, "order_sn = ?", sn).Error
	return
}

// 生成佣金流水
func (srv *OrderCommissionSrv) Do() (err error) {
	done := true
	limit := 100
	offset := 0
	maxCount := 1000
	count := 0
	start, err := util.ChinaYesterday()
	if err != nil {
		return
	}
	// end, err := util.ChinaToday()
	// if err != nil {
	// 	return
	// }
	// 暂时mock为当前时间
	end := time.Now()
	for {
		if count >= maxCount {
			// TODO 输出异常
			break
		}
		args := []interface{}{
			"status = ? AND created_at >= ? AND created_at <= ? AND recommender IS NOT NULL",
			OrderStatusClosed,
			util.FormatTime(start),
			util.FormatTime(end),
		}
		orders, err := orderSrv.List(PGQueryParams{
			Limit:  limit,
			Offset: offset,
			Order:  "createdAt",
		}, args...)
		if err != nil {
			return err
		}
		if len(orders) < limit {
			done = true
		}
		for _, order := range orders {
			e := srv.createOrUpdate(order)
			if e != nil {
				return e
			}
		}

		if done {
			break
		}
		offset += limit
		count++
	}
	return nil
}

// List list order commission
func (srv *OrderCommissionSrv) List(params PGQueryParams, args ...interface{}) (result OrderCommissions, err error) {
	result = make(OrderCommissions, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// Count count order commission
func (srv *OrderCommissionSrv) Count(args ...interface{}) (count int64, err error) {
	return pgCount(&OrderCommission{}, args...)
}

// ListAll list all order commission
func (srv *OrderCommissionSrv) ListAll(params PGQueryParams, args ...interface{}) (result OrderCommissions, err error) {
	result = make(OrderCommissions, 0)
	if params.Limit <= 0 {
		params.Limit = 50
	}
	params.Offset = 0
	params.Order = "-id"
	var tmpResult OrderCommissions
	for {
		tmpResult, err = srv.List(params, args...)
		if err != nil {
			return
		}
		result = append(result, tmpResult...)
		if len(tmpResult) < params.Limit {
			break
		}
		params.Offset += params.Limit
	}
	return
}
