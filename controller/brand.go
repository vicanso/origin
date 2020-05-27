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
	"strconv"
	"strings"

	"github.com/vicanso/elton"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/router"
	"github.com/vicanso/origin/service"
	"github.com/vicanso/origin/util"
	"github.com/vicanso/origin/validate"
)

type (
	brandCtrl struct{}

	addBrandParams struct {
		Name    string `json:"name" validate:"xBrandName"`
		Status  int    `json:"status" validate:"xBrandStatus"`
		Logo    string `json:"logo" validate:"xBrandLogo"`
		Catalog string `json:"catalog" validate:"xBrandCatalog"`
	}
	updateBrandParams struct {
		Name    string `json:"name" validate:"omitempty,xBrandName"`
		Status  int    `json:"status" validate:"omitempty,xBrandStatus"`
		Logo    string `json:"logo" validate:"omitempty,xBrandLogo"`
		Catalog string `json:"catalog" validate:"omitempty,xBrandCatalog"`
	}
	listBrandParams struct {
		Keyword string `json:"keyword" validate:"omitempty,xKeyword"`
		Status  string `json:"status" validate:"omitempty,xBrandStatus"`
		listParams
	}
)

func init() {
	ctrl := brandCtrl{}
	g := router.NewGroup("/brands")
	// 获取品牌状态
	g.GET(
		"/v1/statuses",
		noCacheIfSetNoCache,
		ctrl.listStatuses,
	)

	// 添加品牌
	g.POST(
		"/v1",
		loadUserSession,
		checkMarketingGroup,
		newTracker(cs.ActionBrandAdd),
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
		checkMarketingGroup,
		newTracker(cs.ActionBrandUpdate),
		ctrl.updateByID,
	)
}

func (params *listBrandParams) toConditions() (conditions []interface{}) {
	queryList := make([]string, 0)
	args := make([]interface{}, 0)
	if params.Keyword != "" {
		queryList = append(queryList, "name ILIKE ?")
		args = append(args, "%"+params.Keyword+"%")
	}
	if params.Status != "" {
		queryList = append(queryList, "status = ?")
		args = append(args, params.Status)
	}
	conditions = make([]interface{}, 0)
	if len(queryList) != 0 {
		conditions = append(conditions, strings.Join(queryList, " AND "))
		conditions = append(conditions, args...)
	}
	return
}

func (ctrl brandCtrl) listStatuses(c *elton.Context) (err error) {
	c.CacheMaxAge("5m")
	c.Body = map[string][]*service.BrandStatus{
		"statuses": brandSrv.ListStatuses(),
	}
	return
}

// add add brand
func (ctrl brandCtrl) add(c *elton.Context) (err error) {
	params := addBrandParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	brand := &service.Brand{
		Name:    params.Name,
		Status:  params.Status,
		Logo:    params.Logo,
		Catalog: params.Catalog,
	}

	brand.FirstLetter = util.GetFirstLetter(params.Name)

	err = brandSrv.Add(brand)
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
	c.Body = struct {
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	data, err := brandSrv.FindByID(uint(id))
	if err != nil {
		return
	}
	c.CacheMaxAge("1m")
	c.Body = data
	return
}

// updateByID update brand by id
func (ctrl brandCtrl) updateByID(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	params := updateBrandParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}

	brand := service.Brand{
		Name:    params.Name,
		Status:  params.Status,
		Logo:    params.Logo,
		Catalog: params.Catalog,
	}
	if brand.Name != "" {
		brand.FirstLetter = util.GetFirstLetter(brand.Name)
	}

	err = brandSrv.UpdateByID(uint(id), brand)

	if err != nil {
		return
	}
	c.NoContent()
	return
}
