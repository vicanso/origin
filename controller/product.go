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
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/router"
	"github.com/vicanso/origin/service"
	"github.com/vicanso/origin/validate"
)

type (
	productCtrl struct{}

	addProductParams struct {
		Name       string     `json:"name,omitempty" validate:"xProductName"`
		Price      float64    `json:"price,omitempty" validate:"xProductPrice"`
		Unit       string     `json:"unit,omitempty" validate:"xProductUnit"`
		Catalog    string     `json:"catalog,omitempty" validate:"xProductCatalog"`
		Pics       []string   `json:"pics,omitempty"`
		MainPic    int        `json:"mainPic,omitempty" validate:"omitempty,xProductMainPic"`
		SN         string     `json:"sn,omitempty" validate:"omitempty,xProductSN"`
		Status     int        `json:"status,omitempty" validate:"xStatus"`
		Keywords   string     `json:"keywords,omitempty"`
		Categories []int64    `json:"categories,omitempty"`
		StartedAt  *time.Time `json:"startedAt,omitempty" validate:"required"`
		EndedAt    *time.Time `json:"endedAt,omitempty" validate:"required"`
		// 产地
		Origin string `json:"origin,omitempty" validate:"omitempty,xProductOrigin"`
		// 产品品牌
		Brand uint `json:"brand,omitempty" validate:"omitempty,xProductBrand"`
	}
	updateProductParams struct {
		Name       string     `json:"name,omitempty" validate:"omitempty,xProductName"`
		Price      float64    `json:"price,omitempty" validate:"omitempty,xProductPrice"`
		Unit       string     `json:"unit,omitempty" validate:"omitempty,xProductUnit"`
		Catalog    string     `json:"catalog,omitempty" validate:"omitempty,xProductCatalog"`
		Pics       []string   `json:"pics,omitempty"`
		MainPic    int        `json:"mainPic,omitempty" validate:"omitempty,xProductMainPic"`
		SN         string     `json:"sn,omitempty" validate:"omitempty,xProductSN"`
		Status     int        `json:"status,omitempty" validate:"omitempty,xStatus"`
		Keywords   string     `json:"keywords,omitempty"`
		Categories []int64    `json:"categories,omitempty"`
		StartedAt  *time.Time `json:"startedAt,omitempty"`
		EndedAt    *time.Time `json:"endedAt,omitempty"`
		// 产地
		Origin string `json:"origin,omitempty" validate:"omitempty,xProductOrigin"`
		// 产品品牌
		Brand uint `json:"brand,omitempty" validate:"omitempty,xProductBrand"`
	}
	listProductParams struct {
		listParams
	}
	addProductCategoryParams struct {
		Name    string  `json:"name,omitempty" validate:"xProductCategoryName"`
		Level   int     `json:"level,omitempty" validate:"xProductCategoryLevel"`
		Status  int     `json:"status,omitempty" validate:"xStatus"`
		Belongs []int64 `json:"belongs,omitempty"`
	}
	updateProductCategoryParams struct {
		Name    string  `json:"name,omitempty" validate:"omitempty,xProductCategoryName"`
		Level   int     `json:"level,omitempty" validate:"omitempty,xProductCategoryLevel"`
		Status  int     `json:"status,omitempty" validate:"omitempty,xStatus"`
		Belongs []int64 `json:"belongs,omitempty"`
	}
	listProductCategoryParams struct {
		listParams

		Keyword string `json:"keyword,omitempty" validate:"omitempty,xKeyword"`
		Status  string `json:"status,omitempty" validate:"omitempty,xStatus"`
		Level   string `json:"level,omitempty" validate:"omitempty,xProductCategoryLevel"`
	}
)

func init() {
	ctrl := productCtrl{}
	g := router.NewGroup("/products")
	// 获取产品
	g.GET(
		"/v1",
		noCacheIfSetNoCache,
		ctrl.list,
	)

	// 获取产品分类
	g.GET(
		"/v1/categories",
		noCacheIfSetNoCache,
		ctrl.listCategory,
	)
	// 添加产品分类
	g.POST(
		"/v1/categories",
		loadUserSession,
		checkMarketingGroup,
		newTracker(cs.ActionProductCategoryAdd),
		ctrl.addCategory,
	)
	// 更新产品分类
	g.PATCH(
		"/v1/categories/{id}",
		loadUserSession,
		checkMarketingGroup,
		newTracker(cs.ActionProductCategoryUpdate),
		ctrl.updateCategoryByID,
	)
	// 获取产品分类详情
	g.GET(
		"/v1/categories/{id}",
		noCacheIfSetNoCache,
		ctrl.findCategoryByID,
	)

	// 添加产品
	g.POST(
		"/v1",
		loadUserSession,
		checkMarketingGroup,
		newTracker(cs.ActionProductAdd),
		ctrl.add,
	)
	// 查询产品详情
	g.GET(
		"/v1/{id}",
		noCacheIfSetNoCache,
		ctrl.findByID,
	)
	// 更新产品
	g.PATCH(
		"/v1/{id}",
		loadUserSession,
		checkMarketingGroup,
		newTracker(cs.ActionProductUpdate),
		ctrl.updateByID,
	)
}

func (params listProductParams) toConditions() (conditions []interface{}) {
	return
}

func (params listProductCategoryParams) toConditions() (conditions []interface{}) {
	conds := queryConditions{}
	if params.Keyword != "" {
		conds.add("name ILIKE ?", "%"+params.Keyword+"%")
	}
	if params.Status != "" {
		conds.add("status = ?", params.Status)
	}
	if params.Level != "" {
		conds.add("level = ?", params.Level)
	}
	return conds.toArray()
}

// add add product
func (ctrl productCtrl) add(c *elton.Context) (err error) {
	params := addProductParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	product := &service.Product{
		Name:       params.Name,
		Price:      params.Price,
		Unit:       params.Unit,
		Catalog:    params.Catalog,
		Pics:       params.Pics,
		MainPic:    params.MainPic,
		SN:         params.SN,
		Status:     params.Status,
		Keywords:   params.Keywords,
		Categories: params.Categories,
		StartedAt:  params.StartedAt,
		EndedAt:    params.EndedAt,
		Origin:     params.Origin,
		Brand:      params.Brand,
	}
	err = productSrv.Add(product)
	if err != nil {
		return
	}
	c.Created(product)
	return
}

// list list all product
func (ctrl productCtrl) list(c *elton.Context) (err error) {
	params := listProductParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	count := -1
	args := params.toConditions()
	queryParams := params.toPGQueryParams()
	if queryParams.Offset == 0 {
		count, err = productSrv.Count(args...)
		if err != nil {
			return
		}
	}

	result, err := productSrv.List(queryParams, args...)
	if err != nil {
		return
	}
	c.CacheMaxAge("1m")
	c.Body = struct {
		Products []*service.Product `json:"products"`
		Count    int                `json:"count"`
	}{
		result,
		count,
	}
	return
}

// findByID find product by id
func (ctrl productCtrl) findByID(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	data, err := productSrv.FindByID(uint(id))
	if err != nil {
		return
	}
	c.CacheMaxAge("1m")
	c.Body = data
	return
}

// updateByID update product by id
func (ctrl productCtrl) updateByID(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	params := updateProductParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	product := service.Product{
		Name:       params.Name,
		Price:      params.Price,
		Unit:       params.Unit,
		Catalog:    params.Catalog,
		Pics:       params.Pics,
		MainPic:    params.MainPic,
		SN:         params.SN,
		Status:     params.Status,
		Keywords:   params.Keywords,
		Categories: params.Categories,
		StartedAt:  params.StartedAt,
		EndedAt:    params.EndedAt,
		Origin:     params.Origin,
		Brand:      params.Brand,
	}
	err = productSrv.UpdateByID(uint(id), product)
	if err != nil {
		return
	}
	c.NoContent()
	return
}

func (ctrl productCtrl) addCategory(c *elton.Context) (err error) {
	params := addProductCategoryParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	cat := &service.ProductCategory{
		Name:    params.Name,
		Level:   params.Level,
		Status:  params.Status,
		Belongs: params.Belongs,
	}

	err = productSrv.AddCategory(cat)
	if err != nil {
		return
	}
	c.Created(cat)
	return
}

// listCategory list all category
func (ctrl productCtrl) listCategory(c *elton.Context) (err error) {
	params := listProductCategoryParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	count := -1
	args := params.toConditions()
	queryParams := params.toPGQueryParams()
	if queryParams.Offset == 0 {
		count, err = productSrv.CountCategory(args...)
		if err != nil {
			return
		}
	}
	result, err := productSrv.ListCategory(queryParams, args...)
	if err != nil {
		return
	}

	c.CacheMaxAge("1m")
	c.Body = &struct {
		ProductCategories []*service.ProductCategory `json:"productCategories,omitempty"`
		Count             int                        `json:"count,omitempty"`
	}{
		result,
		count,
	}
	return
}

// updateCategoryByID update category by id
func (ctrl productCtrl) updateCategoryByID(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	params := updateProductCategoryParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	cat := &service.ProductCategory{
		Name:    params.Name,
		Level:   params.Level,
		Status:  params.Status,
		Belongs: params.Belongs,
	}
	err = productSrv.UpdateCategoryByID(uint(id), cat)
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// findCategoryByID find category by id
func (ctrl productCtrl) findCategoryByID(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	data, err := productSrv.FindCategoryByID(uint(id))
	if err != nil {
		return
	}
	c.CacheMaxAge("1m")
	c.Body = data
	return
}
