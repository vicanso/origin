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

import "github.com/vicanso/origin/helper"

const (
	China = iota + 1
)

type (
	// RegionProvince province
	RegionProvince struct {
		helper.Model

		Name    string `json:"name" gorm:"not null;unique_index:idx_province_name"`
		Code    int    `json:"code" gorm:"not null;unique_index:idx_province_code"`
		Country int    `json:"country" gorm:"not null"`
	}
	// RegionCity city
	RegionCity struct {
		helper.Model

		Name     string `json:"name" gorm:"not null"`
		Code     int    `json:"code" gorm:"not null;unique_index:idx_city_code"`
		Province int    `json:"province" gorm:"not null;index:idx_city_province"`
	}
	// RegionArea area
	RegionArea struct {
		helper.Model

		Name string `json:"name" gorm:"not null"`
		Code int    `json:"code" gorm:"not null;unique_index:idx_area_code"`
		City int    `json:"city" gorm:"not null;index:idx_area_city"`
	}

	// RegionStreet street
	RegionStreet struct {
		helper.Model

		Name string `json:"name" gorm:"not null"`
		Code int    `json:"code" gorm:"not null;unique_index:idx_street_code"`
		Area int    `json:"Area" gorm:"not null;index:idx_street_area"`
	}

	RegionSrv struct{}
)

func init() {
	pgGetClient().AutoMigrate(&RegionProvince{}).
		AutoMigrate(&RegionCity{}).
		AutoMigrate(&RegionArea{}).
		AutoMigrate(&RegionStreet{})
}

// AddProvince add province
func (srv *RegionSrv) AddProvince(province *RegionProvince) (err error) {
	err = pgCreate(province)
	return
}

// ListProvince list province
func (srv *RegionSrv) ListProvince(params helper.PGQueryParams, args ...interface{}) (result []*RegionProvince, err error) {
	result = make([]*RegionProvince, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// AddCity add city
func (srv *RegionSrv) AddCity(city *RegionCity) (err error) {
	err = pgCreate(city)
	return
}

// ListCity list city
func (srv *RegionSrv) ListCity(params helper.PGQueryParams, args ...interface{}) (result []*RegionCity, err error) {
	result = make([]*RegionCity, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// AddArea add area
func (srv *RegionSrv) AddArea(area *RegionArea) (err error) {
	err = pgCreate(area)
	return
}

// AddStreet add street
func (srv *RegionSrv) AddStreet(street *RegionStreet) (err error) {
	err = pgCreate(street)
	return
}
