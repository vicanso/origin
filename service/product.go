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
	"time"

	"github.com/lib/pq"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/helper"
)

type (
	Product struct {
		helper.Model

		Name    string         `json:"name" gorm:"type:varchar(30);not null;unique"`
		Price   float64        `json:"price"`
		Unit    string         `json:"unit"`
		Catalog string         `json:"catalog"`
		Pics    pq.StringArray `json:"pics" gorm:"type:text[]"`
		// 主图，从1开始
		MainPic    int            `json:"mainPic"`
		SN         string         `json:"sn"`
		Status     int            `json:"status"`
		Keywords   string         `json:"keywords"`
		Categories pq.StringArray `json:"categories" gorm:"type:text[]"`
		StartedAt  *time.Time     `json:"startedAt"`
		EndedAt    *time.Time     `json:"endedAt"`
		// 产地
		Origin string `json:"origin"`
		// 产品品牌
		Brand uint `json:"brand"`
	}
	ProductStatus struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	ProductSrv struct{}
)

func init() {
	pgGetClient().AutoMigrate(&Product{})
}

func (srv *ProductSrv) createByID(id uint) *Product {
	p := &Product{}
	p.Model.ID = id
	return p
}

// ListStatuses list all product status
func (srv *ProductSrv) ListStatuses() []*ProductStatus {
	return []*ProductStatus{
		&ProductStatus{
			Name:  "启用",
			Value: cs.ProductStatusEnabled,
		},
		&ProductStatus{
			Name:  "禁用",
			Value: cs.ProductStatusDisabled,
		},
	}
}

// Add add product
func (srv *ProductSrv) Add(product *Product) (err error) {
	err = pgCreate(product)
	return
}

// UpdateByID update product by id
func (srv *ProductSrv) UpdateByID(id uint, attrs ...interface{}) (err error) {
	err = pgGetClient().Model(srv.createByID(id)).Update(attrs...).Error
	return
}

// FindByID find product by id
func (srv *ProductSrv) FindByID(id uint) (product *Product, err error) {
	product = new(Product)
	err = pgGetClient().First(product, "id = ?", id).Error
	return
}

// List list product
func (srv *ProductSrv) List(params helper.PGQueryParams, args ...interface{}) (result []*Product, err error) {
	result = make([]*Product, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// Count count the product
func (srv *ProductSrv) Count(args ...interface{}) (count int, err error) {
	return pgCount(&Product{}, args...)
}
