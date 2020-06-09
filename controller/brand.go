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
	brandCtrl struct{}

	addBrandParams struct {
		Name    string `json:"name,omitempty" validate:"xBrandName"`
		Status  int    `json:"status,omitempty" validate:"xStatus"`
		Logo    string `json:"logo,omitempty" validate:"omitempty,xFile"`
		Catalog string `json:"catalog,omitempty" validate:"xBrandCatalog"`
	}
	updateBrandParams struct {
		Name    string `json:"name,omitempty" validate:"omitempty,xBrandName"`
		Status  int    `json:"status,omitempty" validate:"omitempty,xStatus"`
		Logo    string `json:"logo,omitempty" validate:"omitempty,xFile"`
		Catalog string `json:"catalog,omitempty" validate:"omitempty,xBrandCatalog"`
	}
	listBrandParams struct {
		listParams

		Keyword string `json:"keyword,omitempty" validate:"omitempty,xKeyword"`
		Status  string `json:"status,omitempty" validate:"omitempty,xStatus"`
	}
)

func init() {
	ctrl := brandCtrl{}
	g := router.NewGroup("/brands")
	// 添加品牌
	g.POST(
		"/v1",
		loadUserSession,
		newTracker(cs.ActionBrandAdd),
		checkMarketingGroup,
		ctrl.add,
	)
	// 品牌列表
	g.GET(
		"/v1",
		noCacheIfSetNoCache,
		ctrl.list,
	)
	// 获取品牌详细信息
	g.GET(
		"/v1/{id}",
		noCacheIfSetNoCache,
		ctrl.findByID,
	)
	// 更新品牌信息
	g.PATCH(
		"/v1/{id}",
		loadUserSession,
		newTracker(cs.ActionBrandUpdate),
		checkMarketingGroup,
		ctrl.updateByID,
	)
}

func (params listBrandParams) toConditions() []interface{} {
	conds := queryConditions{}

	if params.Keyword != "" {
		conds.add("name ILIKE ?", "%"+params.Keyword+"%")
	}
	if params.Status != "" {
		conds.add("status = ?", params.Status)
	}
	return conds.toArray()
}

// add add brand
func (ctrl brandCtrl) add(c *elton.Context) (err error) {
	params := addBrandParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}

	brand, err := brandSrv.Add(service.Brand{
		Name:    params.Name,
		Status:  params.Status,
		Logo:    params.Logo,
		Catalog: params.Catalog,
	})
	if err != nil {
		return
	}
	c.Created(brand)

	return
}

// list list all brand
func (ctrl brandCtrl) list(c *elton.Context) (err error) {
	params := listBrandParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	count := -1
	args := params.toConditions()
	queryParams := params.toPGQueryParams()
	if queryParams.Offset == 0 {
		count, err = brandSrv.Count(args...)
		if err != nil {
			return
		}
	}
	result, err := brandSrv.List(queryParams, args...)
	if err != nil {
		return
	}
	c.CacheMaxAge("1m")
	c.Body = &struct {
		Brands []*service.Brand `json:"brands"`
		Count  int              `json:"count"`
	}{
		result,
		count,
	}

	return
}

// findByID find brand by id
func (ctrl brandCtrl) findByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	data, err := brandSrv.FindByID(id)
	if err != nil {
		return
	}
	c.CacheMaxAge("1m")
	c.Body = data
	return
}

// updateByID update brand by id
func (ctrl brandCtrl) updateByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	params := updateBrandParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}

	err = brandSrv.UpdateByID(id, service.Brand{
		Name:    params.Name,
		Status:  params.Status,
		Logo:    params.Logo,
		Catalog: params.Catalog,
	})

	if err != nil {
		return
	}
	c.NoContent()
	return
}
