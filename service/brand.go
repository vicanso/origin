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

	lruTTL "github.com/vicanso/lru-ttl"
	"github.com/vicanso/origin/helper"
	"github.com/vicanso/origin/util"
	"gorm.io/gorm"
)

type (
	Brands []*Brand
	Brand  struct {
		helper.Model

		Name        string `json:"name,omitempty" gorm:"type:varchar(100);not null;unique_index:idx_brand_name"`
		Status      int    `json:"status,omitempty" gorm:"index:idx_brand_status"`
		StatusDesc  string `json:"statusDesc,omitempty" gorm:"-"`
		Logo        string `json:"logo,omitempty"`
		Catalog     string `json:"catalog,omitempty"`
		FirstLetter string `json:"firstLetter,omitempty"`
	}
	BrandSrv struct{}
)

var (
	// brandNameCache brandh's name cache
	brandNameCache *lruTTL.Cache
)

func init() {
	ttl := 300 * time.Second
	// 本地开发环境，设置缓存为1秒
	if util.IsDevelopment() {
		ttl = time.Second
	}
	brandNameCache = lruTTL.New(500, ttl)
	err := helper.PGAutoMigrate(&Brand{})
	if err != nil {
		panic(err)
	}
}

func (bs Brands) AfterFind(tx *gorm.DB) (err error) {
	for _, b := range bs {
		err = b.AfterFind(tx)
		if err != nil {
			return
		}
	}
	return
}

func (b *Brand) AfterFind(_ *gorm.DB) (err error) {
	b.StatusDesc = getStatusDesc(b.Status)
	return
}

func (b *Brand) BeforeCreate(_ *gorm.DB) (err error) {
	// 自动生成拼音首字母
	if b.Name != "" {
		b.FirstLetter = util.GetFirstLetter(b.Name)
	}
	return
}
func (b *Brand) AfterUpdate(tx *gorm.DB) (err error) {
	if b.Name != "" {
		letter := util.GetFirstLetter(b.Name)
		// 自动更新拼音首字母
		if b.FirstLetter != letter {
			tx.Model(b).Updates(Brand{
				FirstLetter: letter,
			})
		}
	}
	return
}

func (srv *BrandSrv) createByID(id uint) *Brand {
	b := &Brand{}
	b.Model.ID = id
	return b
}

// Add add brand
func (srv *BrandSrv) Add(data Brand) (brand *Brand, err error) {
	brand = &data
	err = pgCreate(brand)
	return
}

// UpdateByID update brand by id
func (srv *BrandSrv) UpdateByID(id uint, brand Brand) (err error) {
	err = pgGetClient().Model(srv.createByID(id)).Updates(brand).Error
	return
}

// FindByID find brand by id
func (srv *BrandSrv) FindByID(id uint) (brand *Brand, err error) {
	brand = new(Brand)
	err = pgGetClient().First(brand, "id = ?", id).Error
	return
}

// List list brands
func (srv *BrandSrv) List(params PGQueryParams, args ...interface{}) (result Brands, err error) {
	result = make(Brands, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// Count count the brand
func (srv *BrandSrv) Count(args ...interface{}) (count int64, err error) {
	return pgCount(&Brand{}, args...)
}

// GetNameFromCache get name from cache
// If not exists, it willl get from db and save it to cache
func (srv *BrandSrv) GetNameFromCache(id uint) (name string, err error) {
	if id == 0 {
		return
	}
	value, ok := brandNameCache.Get(id)
	if ok {
		return value.(string), nil
	}
	brand, err := srv.FindByID(id)
	if err != nil {
		return
	}
	name = brand.Name
	brandNameCache.Add(id, name)
	return
}
