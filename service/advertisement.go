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

	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/helper"
	"gorm.io/gorm"
)

var (
	// 广告类别
	advertisementCategoriesMap map[string]string
)

type (
	// Advertisement 广告
	Advertisement struct {
		helper.Model

		Status     int    `json:"status,omitempty" gorm:"index:idx_advertisement_status"`
		StatusDesc string `json:"statusDesc,omitempty" gorm:"-"`
		// 链接
		Link string `json:"link,omitempty" gorm:"not null"`
		// 描述
		Summary      string `json:"summary,omitempty"`
		Category     string `json:"category,omitempty" gorm:"index:idx_advertisement_category;not null"`
		CategoryDesc string `json:"categoryDesc,omitempty" gorm:"-"`
		// 排序
		Rank int `json:"rank,omitempty"`

		// 图片
		Pic string `json:"pic,omitempty" gorm:"not null"`

		// 开始结束时间
		StartedAt *time.Time `json:"startedAt,omitempty" gorm:"not null"`
		EndedAt   *time.Time `json:"endedAt,omitempty" gorm:"not null"`
	}
	// AdvertisementCategory 广告分类
	AdvertisementCategory struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}
	AdvertisementSrv struct{}
)

func init() {
	advertisementCategoriesMap = map[string]string{
		cs.AdvertisementHome:     "首页",
		cs.AdvertisementCategory: "分类页",
	}
	err := pgGetClient().AutoMigrate(&Advertisement{})
	if err != nil {
		panic(err)
	}
}

func (a *Advertisement) AfterFind(_ *gorm.DB) (err error) {
	a.StatusDesc = getStatusDesc(a.Status)
	a.CategoryDesc = advertisementCategoriesMap[a.Category]
	return
}

func (srv *AdvertisementSrv) createByID(id uint) *Advertisement {
	a := &Advertisement{}
	a.Model.ID = id
	return a
}

// ListCategories list advertisement category
func (srv *AdvertisementSrv) ListCategory() []*AdvertisementCategory {
	advertisementCategories := make([]*AdvertisementCategory, 0)
	for key, value := range advertisementCategoriesMap {
		advertisementCategories = append(advertisementCategories, &AdvertisementCategory{
			Name:  value,
			Value: key,
		})
	}
	return advertisementCategories
}

// Add add advertisement
func (srv *AdvertisementSrv) Add(data Advertisement) (advertisement *Advertisement, err error) {
	advertisement = &data
	err = pgCreate(advertisement)
	return
}

// UpdateByID update advertisement
func (srv *AdvertisementSrv) UpdateByID(id uint, advertisement Advertisement) (err error) {
	err = pgGetClient().Model(srv.createByID(id)).Updates(advertisement).Error
	return
}

// FindByID find advertisement by id
func (srv *AdvertisementSrv) FindByID(id uint) (advertisement *Advertisement, err error) {
	advertisement = new(Advertisement)
	err = pgGetClient().First(advertisement, "id = ?", id).Error
	return
}

// List list advertisement
func (srv *AdvertisementSrv) List(params PGQueryParams, args ...interface{}) (result []*Advertisement, err error) {
	result = make([]*Advertisement, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// Count count the advertisement
func (srv *AdvertisementSrv) Count(args ...interface{}) (count int64, err error) {
	return pgCount(&Advertisement{}, args...)
}
