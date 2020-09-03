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
	"github.com/vicanso/origin/helper"
	"gorm.io/gorm"
)

type (
	Suppliers []*Supplier
	Supplier  struct {
		helper.Model

		Name            string `json:"name,omitempty" gorm:"type:varchar(50);not null;index:idx_supplier_name"`
		BaseAddress     string `json:"baseAddress,omitempty"`
		BaseAddressDesc string `json:"baseAddressDesc,omitempty" gorm:"-"`
		Address         string `json:"address,omitempty"`
		Mobile          string `json:"mobile,omitempty"`
		Contact         string `json:"contact,omitempty"`

		// 状态
		Status int `json:"status,omitempty" gorm:"index:idx_supplier_status"`
		// 状态描述
		StatusDesc string `json:"statusDesc,omitempty" gorm:"-"`
	}
	SupplierSrv struct{}
)

func init() {
	err := helper.PGAutoMigrate(&Supplier{})
	if err != nil {
		panic(err)
	}
}

func (s *Supplier) AfterFind(_ *gorm.DB) (err error) {
	s.StatusDesc = getStatusDesc(s.Status)
	s.BaseAddressDesc, _ = regionSrv.GetNameFromCache(s.BaseAddress, 0)
	return
}

func (suppliers Suppliers) AfterFind(tx *gorm.DB) (err error) {
	for _, s := range suppliers {
		err = s.AfterFind(tx)
		if err != nil {
			return
		}
	}
	return
}

func (srv *SupplierSrv) createByID(id uint) *Supplier {
	s := &Supplier{}
	s.Model.ID = id
	return s
}

// Add add supplier
func (srv *SupplierSrv) Add(data Supplier) (supplier *Supplier, err error) {
	supplier = &data
	err = pgCreate(supplier)
	return
}

// UpdateByID update supplier by id
func (srv *SupplierSrv) UpdateByID(id uint, supplier Supplier) (err error) {
	err = pgGetClient().Model(srv.createByID(id)).Updates(supplier).Error
	return
}

// FindByID find supplier by id
func (srv *SupplierSrv) FindByID(id uint) (supplier *Supplier, err error) {
	supplier = new(Supplier)
	err = pgGetClient().First(supplier, "id = ?", id).Error
	return
}

// List list supplier
func (srv *SupplierSrv) List(params PGQueryParams, args ...interface{}) (result []*Supplier, err error) {
	result = make([]*Supplier, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// Count count the supplier
func (srv *SupplierSrv) Count(args ...interface{}) (count int64, err error) {
	return pgCount(&Supplier{}, args...)
}
