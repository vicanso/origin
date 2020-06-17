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
	"strings"
	"time"

	lruTTL "github.com/vicanso/lru-ttl"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/helper"
	"github.com/vicanso/origin/util"
)

type (
	Region struct {
		helper.Model

		Category   int    `json:"category,omitempty" gorm:"not null;index:idx_region_category_name"`
		Name       string `json:"name,omitempty" gorm:"not null;index:idx_region_category_name"`
		Code       string `json:"code,omitempty" gorm:"not null;unique_index:idx_region_code"`
		Parent     string `json:"parent,omitempty" gorm:"not null;index:idx_region_parent"`
		Status     int    `json:"status,omitempty"`
		StatusDesc string `json:"statusDesc,omitempty" gorm:"-"`
		Priority   int    `json:"priority,omitempty"`
	}
	RegionCategory struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}

	RegionSrv struct{}
)

var (
	// regionNameCache region's cache
	regionNameCache *lruTTL.Cache
)

func init() {
	ttl := 300 * time.Second
	// 本地开发环境，设置缓存为1秒
	if util.IsDevelopment() {
		ttl = time.Second
	}
	regionNameCache = lruTTL.New(200, ttl)
	pgGetClient().AutoMigrate(&Region{})
}

func (r *Region) AfterFind() (err error) {
	r.StatusDesc = getStatusDesc(r.Status)
	return
}

func (r *Region) BeforeCreate() (err error) {
	if r.Priority == 0 {
		r.Priority = 1
	}
	return
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
func (srv *RegionSrv) Add(data Region) (region *Region, err error) {
	region = &data
	err = pgCreate(region)
	return
}

// List list region
func (srv *RegionSrv) List(params PGQueryParams, args ...interface{}) (result []*Region, err error) {
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

// FindByCode find by code
func (srv *RegionSrv) FindByCode(code string) (region *Region, err error) {
	region = new(Region)
	err = pgGetClient().First(region, "code = ?", code).Error
	return
}

// UpdateByID update region by id
func (srv *RegionSrv) UpdateByID(id uint, region Region) (err error) {
	err = pgGetClient().Model(srv.createByID(id)).Updates(region).Error
	return
}

// GetNameFromCache get region's name from cache
// If not exists, it will get from db and save to cache
func (srv *RegionSrv) GetNameFromCache(code string, startLevel int) (name string, err error) {
	if code == "" {
		return
	}
	value, ok := regionNameCache.Get(code)
	if ok {
		return value.(string), nil
	}
	// 最多查询5层
	regions := make([]string, 0)
	searchCode := code
	for i := 0; i < 5; i++ {
		region, e := srv.FindByCode(searchCode)
		if e != nil {
			err = e
			return
		}
		regions = append([]string{
			region.Name,
		}, regions...)
		searchCode = region.Parent
		if searchCode == "" {
			break
		}
	}
	name = strings.Join(regions[startLevel:], "")
	regionNameCache.Add(code, name)
	return
}
