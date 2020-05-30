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
	Region struct {
		helper.Model

		Category int    `json:"category,omitempty" gorm:"not null;index:idx_region_category_name"`
		Name     string `json:"name,omitempty" gorm:"not null;index:idx_region_category_name"`
		Code     string `json:"code,omitempty" gorm:"not null;unique_index:idx_region_code"`
		Parent   string `json:"parent,omitempty" gorm:"not null;index:idx_region_parent"`
		Status   int    `json:"status,omitempty"`
	}
	RegionStatus struct {
		Name  string `json:"name,omitempty"`
		Value int    `json:"value,omitempty"`
	}
	RegionCategory struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}

	RegionSrv struct{}
)

func init() {
	pgGetClient().AutoMigrate(&Region{})
}

// GetCategoryIndex get index of category
func (srv *RegionSrv) GetCategoryIndex(category string) (index int) {
	for i, item := range cs.RegionCategories {
		if item == category {
			index = i + 1
		}
	}
	return
}

// ListStatus list status
func (srv *RegionSrv) ListStatus() []*RegionStatus {
	return []*RegionStatus{
		{
			Name:  "启用",
			Value: cs.RegionStatusEnabled,
		},
		{
			Name:  "禁用",
			Value: cs.RegionStatusDisabled,
		},
	}
}

// ListCategory list category
func (srv *RegionSrv) ListCategory() []*RegionCategory {
	return []*RegionCategory{
		{
			Name:  "国",
			Value: cs.RegionCountry,
		},
		{
			Name:  "省",
			Value: cs.RegionProvince,
		},
		{
			Name:  "市",
			Value: cs.RegionCity,
		},
		{
			Name:  "区",
			Value: cs.RegionArea,
		},
		{
			Name:  "街",
			Value: cs.RegionStreet,
		},
	}
}

func (srv *RegionSrv) createByID(id uint) *Region {
	r := &Region{}
	r.Model.ID = id
	return r
}

// Add add region
func (srv *RegionSrv) Add(region *Region) (err error) {
	err = pgCreate(region)
	return
}

// List list region
func (srv *RegionSrv) List(params helper.PGQueryParams, args ...interface{}) (result []*Region, err error) {
	result = make([]*Region, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// Count count region
func (srv *RegionSrv) Count(args ...interface{}) (count int, err error) {
	return pgCount(&Region{}, args...)
}

// FindByID find by id
func (srv *RegionSrv) FindByID(id uint) (region *Region, err error) {
	region = new(Region)
	err = pgGetClient().First(region, "id = ?", id).Error
	return
}

// UpdateByID update region by id
func (srv *RegionSrv) UpdateByID(id uint, attrs ...interface{}) (err error) {
	err = pgGetClient().Model(srv.createByID(id)).Update(attrs...).Error
	return
}
