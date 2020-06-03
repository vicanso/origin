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
	"github.com/vicanso/elton"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/router"
	"github.com/vicanso/origin/service"
	"github.com/vicanso/origin/validate"
)

type (
	supplierCtrl struct{}

	addSupplierParams struct {
		Name    string `json:"name,omitempty" validate:"xSupplierName"`
		Address string `json:"address,omitempty" validate:"xSupplierAddress"`
		Mobile  string `json:"mobile,omitempty" validate:"xMobile"`
		Contact string `json:"contact,omitempty" validate:"xSupplierContact"`
		Status  int    `json:"status,omitempty" validate:"xStatus"`
	}
	updateSupplierParms struct {
		Name    string `json:"name,omitempty" validate:"omitempty,xSupplierName"`
		Address string `json:"address,omitempty" validate:"omitempty,xSupplierAddress"`
		Mobile  string `json:"mobile,omitempty" validate:"omitempty,xMobile"`
		Contact string `json:"contact,omitempty" validate:"omitempty,xSupplierContact"`
		Status  int    `json:"status,omitempty" validate:"omitempty,xStatus"`
	}
	listSupplierParams struct {
		listParams

		Keyword string `json:"keyword,omitempty" validate:"omitempty,xKeyword"`
		Status  string `json:"status,omitempty" validate:"omitempty,xStatus"`
	}
)

func init() {
	g := router.NewGroup("/suppliers")
	ctrl := supplierCtrl{}

	g.GET(
		"/v1",
		loadUserSession,
		checkMarketingGroup,
		ctrl.list,
	)
	g.POST(
		"/v1",
		loadUserSession,
		checkMarketingGroup,
		newTracker(cs.ActionSupplierAdd),
		ctrl.add,
	)

	g.PATCH(
		"/v1/{id}",
		loadUserSession,
		checkMarketingGroup,
		newTracker(cs.ActionSupplierUpdate),
		ctrl.updateByID,
	)
	g.GET(
		"/v1/{id}",
		loadUserSession,
		checkMarketingGroup,
		ctrl.findByID,
	)
}

func (params listSupplierParams) toConditions() []interface{} {
	conds := queryConditions{}

	if params.Keyword != "" {
		conds.add("name ILIKE ?", "%"+params.Keyword+"%")
	}
	if params.Status != "" {
		conds.add("status = ?", params.Status)
	}
	return conds.toArray()
}

// add add supplier
func (supplierCtrl) add(c *elton.Context) (err error) {
	params := addSupplierParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	supplier, err := supplierSrv.Add(service.Supplier{
		Name:    params.Name,
		Address: params.Address,
		Mobile:  params.Mobile,
		Contact: params.Contact,
		Status:  params.Status,
	})
	if err != nil {
		return
	}
	c.Created(supplier)
	return
}

// updateByID update supplier by id
func (supplierCtrl) updateByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	params := updateSupplierParms{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	err = supplierSrv.UpdateByID(id, &service.Supplier{
		Name:    params.Name,
		Address: params.Address,
		Mobile:  params.Mobile,
		Contact: params.Contact,
		Status:  params.Status,
	})
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// findByID find supplier by id
func (supplierCtrl) findByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	supplier, err := supplierSrv.FindByID(id)
	if err != nil {
		return
	}
	c.Body = supplier
	return
}

// list list supplier
func (supplierCtrl) list(c *elton.Context) (err error) {
	params := listSupplierParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	count := -1
	args := params.toConditions()
	queryParams := params.toPGQueryParams()
	if queryParams.Offset == 0 {
		count, err = supplierSrv.Count(args...)
		if err != nil {
			return
		}
	}
	result, err := supplierSrv.List(queryParams, args...)
	if err != nil {
		return
	}
	c.Body = &struct {
		Suppliers []*service.Supplier `json:"suppliers,omitempty"`
		Count     int                 `json:"count,omitempty"`
	}{
		result,
		count,
	}

	return
}
