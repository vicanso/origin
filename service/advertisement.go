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

	"github.com/vicanso/origin/helper"
)

type (
	Advertisement struct {
		helper.Model

		Status     int    `json:"status,omitempty" gorm:"index:idx_advertisement_status"`
		StatusDesc string `json:"statusDesc,omitempty" gorm:"-"`
		// 链接
		Link string `json:"link,omitempty" gorm:"not null"`
		// 描述
		Summary  string `json:"summary,omitempty"`
		Category string `json:"category,omitempty" gorm:"index:idx_advertisement_category;not null"`

		// 开始结束时间
		StartedAt *time.Time `json:"startedAt,omitempty" gorm:"not null"`
		EndedAt   *time.Time `json:"endedAt,omitempty" gorm:"not null"`
	}
	AdvertisementSrv struct{}
)

func init() {
	pgGetClient().AutoMigrate(&Advertisement{})
}

func (a *Advertisement) AfterFind() (err error) {
	a.StatusDesc = getStatusDesc(a.Status)
	return
}

func (srv *AdvertisementSrv) createByID(id uint) *Advertisement {
	a := &Advertisement{}
	a.Model.ID = id
	return a
}

// Add add advertisement
func (srv *AdvertisementSrv) Add(data Advertisement) (advertisement *Advertisement, err error) {
	advertisement = &data
	err = pgCreate(advertisement)
	return
}

// UpdateByID update advertisement
func (srv *AdvertisementSrv) UpdateByID(id uint, attrs ...interface{}) (err error) {
	err = pgGetClient().Model(srv.createByID(id)).Update(attrs...).Error
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
func (srv *AdvertisementSrv) Count(args ...interface{}) (count int, err error) {
	return pgCount(&Advertisement{}, args...)
}
