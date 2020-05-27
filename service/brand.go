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
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/helper"
)

type (
	Brand struct {
		helper.Model

		Name        string `json:"name" gorm:"type:varchar(100);not null;unique_index:idx_brand_name"`
		Status      int    `json:"status"`
		Logo        string `json:"logo"`
		Catalog     string `json:"catalog"`
		FirstLetter string `json:"firstLetter"`
	}
	BrandStatus struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	BrandSrv struct{}
)

func init() {
	pgGetClient().AutoMigrate(&Brand{})
}

func (srv *BrandSrv) createByID(id uint) *Brand {
	b := &Brand{}
	b.Model.ID = id
	return b
}

// ListStatuses list all brand status
func (srv *BrandSrv) ListStatuses() []*BrandStatus {
	return []*BrandStatus{
		&BrandStatus{
			Name:  "启用",
			Value: cs.BrandStatusEnabled,
		},
		&BrandStatus{
			Name:  "禁用",
			Value: cs.BrandStatusDisabled,
		},
	}
}

// Add add brand
func (srv *BrandSrv) Add(brand *Brand) (err error) {
	err = pgCreate(brand)
	return
}

// UpdateByID update brand by id
func (srv *BrandSrv) UpdateByID(id uint, attrs ...interface{}) (err error) {
	err = pgGetClient().Model(srv.createByID(id)).Update(attrs...).Error
	return
}

// FindByID find brand by id
func (srv *BrandSrv) FindByID(id uint) (brand *Brand, err error) {
	brand = new(Brand)
	err = pgGetClient().First(brand, "id = ?", id).Error
	return
}

// List list brands
func (srv *BrandSrv) List(params helper.PGQueryParams, args ...interface{}) (result []*Brand, err error) {
	result = make([]*Brand, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// Count count the brand
func (srv *BrandSrv) Count(args ...interface{}) (count int, err error) {
	return pgCount(&Brand{}, args...)
}
