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

type (
	Receiver struct {
		helper.Model

		UserID          uint   `json:"userID,omitempty" gorm:"index:idx_receiver_user;not null"`
		Name            string `json:"name,omitempty"`
		Mobile          string `json:"mobile,omitempty"`
		BaseAddress     string `json:"baseAddress,omitempty"`
		BaseAddressDesc string `json:"baseAddressDesc,omitempty" gorm:"-"`
		Address         string `json:"address,omitempty"`
	}
	ReceiverSrv struct{}
)

func init() {
	pgGetClient().AutoMigrate(&Receiver{})
}

func (r *Receiver) AfterFind() (err error) {
	r.BaseAddressDesc, _ = regionSrv.GetNameFromCache(r.BaseAddress, 1)
	return
}

func (srv *ReceiverSrv) createByID(id uint) *Receiver {
	r := Receiver{}
	r.Model.ID = id
	return &r
}

// Add add receiver
func (srv *ReceiverSrv) Add(data Receiver) (receiver *Receiver, err error) {
	receiver = &data
	err = pgCreate(receiver)
	return
}

// UpdateByID update receiver by id
func (srv *ReceiverSrv) UpdateByID(id uint, attrs ...interface{}) (err error) {
	err = pgGetClient().Model(srv.createByID(id)).Update(attrs...).Error
	return
}

// FindByID find receiver by id
func (srv *ReceiverSrv) FindByID(id uint) (receiver *Receiver, err error) {
	receiver = new(Receiver)
	err = pgGetClient().First(receiver, "id = ?", id).Error
	return
}

// List list receiver
func (srv *ReceiverSrv) List(params PGQueryParams, args ...interface{}) (result []*Receiver, err error) {
	result = make([]*Receiver, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// Count count the receiver
func (srv *ReceiverSrv) Count(args ...interface{}) (count int, err error) {
	return pgCount(&Receiver{}, args...)
}
