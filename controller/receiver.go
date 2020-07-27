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

package controller

import (
	"net/http"

	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/router"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	"github.com/vicanso/origin/service"
	"github.com/vicanso/origin/validate"
)

type (
	receiverCtrl struct{}

	addReceiverParams struct {
		Name        string `json:"name,omitempty" validate:"xReceiverName"`
		Mobile      string `json:"mobile,omitempty" validate:"xMobile"`
		BaseAddress string `json:"baseAddress,omitempty" validate:"xBaseAddress"`
		Address     string `json:"address,omitempty" validate:"xAddress"`
	}
	updateReceiverParams struct {
		Name        string `json:"name,omitempty" validate:"omitempty,xReceiverName"`
		Mobile      string `json:"mobile,omitempty" validate:"omitempty,xMobile"`
		BaseAddress string `json:"baseAddress,omitempty" validate:"omitempty,xBaseAddress"`
		Address     string `json:"address,omitempty" validate:"omitempty,xAddress"`
	}
)

const (
	errReceiverCategory = "receiver"
)

var (
	errReceiverUserInvalid = &hes.Error{
		Message:    "该账户不能修改收货人信息",
		StatusCode: http.StatusBadRequest,
		Category:   errReceiverCategory,
	}
)

func init() {
	g := router.NewGroup("/receivers", loadUserSession, shouldBeLogined)
	ctrl := receiverCtrl{}

	g.GET(
		"/v1",
		ctrl.list,
	)

	g.POST(
		"/v1",
		newTracker(cs.ActionReceiverAdd),
		ctrl.add,
	)
	g.PATCH(
		"/v1/{id}",
		newTracker(cs.ActionReceiverUpdate),
		ctrl.updateByID,
	)
	g.DELETE(
		"/v1/{id}",
		newTracker(cs.ActionConfigurationDelete),
		ctrl.deleteByID,
	)

}

// add add receiver
func (receiverCtrl) add(c *elton.Context) (err error) {
	params := addReceiverParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	us := getUserSession(c)
	// TODO 是否限制可添加的最多收货人数量
	receiver, err := receiverSrv.Add(service.Receiver{
		UserID:      us.GetID(),
		Name:        params.Name,
		Mobile:      params.Mobile,
		BaseAddress: params.BaseAddress,
		Address:     params.Address,
	})
	if err != nil {
		return
	}
	c.Created(receiver)
	return
}

func (receiverCtrl) validateOwner(c *elton.Context, id uint) (err error) {
	result, err := receiverSrv.FindByID(id)
	if err != nil {
		return
	}
	us := getUserSession(c)
	if result.UserID != us.GetID() {
		err = errReceiverUserInvalid
		return
	}
	return
}

// updateByID update receiver by id
func (ctrl receiverCtrl) updateByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	params := updateReceiverParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	err = ctrl.validateOwner(c, id)
	if err != nil {
		return
	}
	err = receiverSrv.UpdateByID(id, service.Receiver{
		Name:        params.Name,
		Mobile:      params.Mobile,
		BaseAddress: params.BaseAddress,
		Address:     params.Address,
	})
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// list list receiver
func (receiverCtrl) list(c *elton.Context) (err error) {
	us := getUserSession(c)
	conds := queryConditions{}
	conds.add("user_id = ?", us.GetID())
	result, err := receiverSrv.List(PGQueryParams{}, conds.toArray()...)
	if err != nil {
		return
	}
	c.Body = &struct {
		Receivers service.Receivers `json:"receivers,omitempty"`
	}{
		result,
	}
	return
}

// delete delete receiver
func (ctrl receiverCtrl) deleteByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	err = ctrl.validateOwner(c, id)
	if err != nil {
		return
	}
	err = receiverSrv.Delete(id)
	if err != nil {
		return
	}
	c.NoContent()
	return
}
