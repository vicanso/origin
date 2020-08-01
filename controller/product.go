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
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/helper"
	"github.com/vicanso/origin/router"
	"github.com/vicanso/origin/service"
	"github.com/vicanso/origin/validate"
)

type (
	productCtrl struct{}

	addProductParams struct {
		Name       string     `json:"name,omitempty" validate:"xProductName"`
		Price      float64    `json:"price,omitempty" validate:"xProductPrice"`
		Specs      uint       `json:"specs,omitempty" validate:"xProductSpecs"`
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
		// 供应商
		Supplier uint `json:"supplier,omitempty" validate:"xProductSupplier"`
		// 排序
		Rank int `json:"rank,omitempty" validate:"omitempty,xRank"`
	}
	updateProductParams struct {
		Name       string     `json:"name,omitempty" validate:"omitempty,xProductName"`
		Price      float64    `json:"price,omitempty" validate:"omitempty,xProductPrice"`
		Specs      uint       `json:"specs,omitempty" validate:"omitempty,xProductSpecs"`
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
		// 供应商
		Supplier uint `json:"supplier,omitempty" validate:"omitempty,xProductSupplier"`
		// 排序
		Rank int `json:"rank,omitempty" validate:"omitempty,xRank"`
	}
	listProductParams struct {
		listParams

		IDS         string `json:"ids,omitempty"`
		Category    string `json:"category,omitempty" validate:"omitempty,xProductCategory"`
		Status      string `json:"status,omitempty" validate:"omitempty,xStatus"`
		Purchasable string `json:"purchasable,omitempty"`
		Keyword     string `json:"keyword,omitempty" validate:"omitempty,xKeyword"`
	}
	addProductCategoryParams struct {
		Name    string  `json:"name,omitempty" validate:"xProductCategoryName"`
		Level   int     `json:"level,omitempty" validate:"xProductCategoryLevel"`
		Status  int     `json:"status,omitempty" validate:"xStatus"`
		Belongs []int64 `json:"belongs,omitempty"`
		Rank    int     `json:"rank,omitempty" validate:"omitempty,xRank"`
		Icon    string  `json:"icon,omitempty" validate:"xFile"`
	}
	updateProductCategoryParams struct {
		Name    string  `json:"name,omitempty" validate:"omitempty,xProductCategoryName"`
		Level   int     `json:"level,omitempty" validate:"omitempty,xProductCategoryLevel"`
		Status  int     `json:"status,omitempty" validate:"omitempty,xStatus"`
		Belongs []int64 `json:"belongs,omitempty"`
		Rank    int     `json:"rank,omitempty" validate:"omitempty,xRank"`
		Icon    string  `json:"icon,omitempty" validate:"omitempty,xFile"`
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

	// 获取热门产品搜索关键字
	g.GET(
		"/v1/search-hot-keywords",
		ctrl.listSearchHotKeywords,
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
		newTracker(cs.ActionProductCategoryAdd),
		checkMarketingGroup,
		ctrl.addCategory,
	)
	// 更新产品分类
	g.PATCH(
		"/v1/categories/{id}",
		loadUserSession,
		newTracker(cs.ActionProductCategoryUpdate),
		checkMarketingGroup,
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
		newTracker(cs.ActionProductAdd),
		checkMarketingGroup,
		ctrl.add,
	)
	// 查询产品详情
	g.GET(
		"/v1/{id}",
		noCacheIfSetNoCache,
		ctrl.findByID,
	)
	// 获取产品主图
	g.GET(
		"/v1/{id}/cover/{quality}-{width}-{height}.{ext}",
		ctrl.getMainImage,
	)
	// 更新产品
	g.PATCH(
		"/v1/{id}",
		loadUserSession,
		newTracker(cs.ActionProductUpdate),
		checkMarketingGroup,
		ctrl.updateByID,
	)
}

func (params listProductParams) toConditions() (conditions []interface{}) {
	conds := queryConditions{}
	if params.IDS != "" {
		conds.add("id IN (?)", strings.Split(params.IDS, ","))
	}
	if params.Category != "" {
		conds.add("? = ANY(categories)", params.Category)
	}
	if params.Status != "" {
		conds.add("status = ?", params.Status)
	}
	if params.Purchasable != "" {
		now := time.Now()
		conds.add("status = ?", cs.StatusEnabled)
		conds.add("started_at < ?", now)
		conds.add("ended_at > ?", now)
	}

	if params.Keyword != "" {
		conds.add("name ILIKE ?", "%"+params.Keyword+"%")
	}

	return conds.toArray()
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
	product, err := productSrv.Add(service.Product{
		Name:       params.Name,
		Price:      params.Price,
		Specs:      params.Specs,
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
		Supplier:   params.Supplier,
		Rank:       params.Rank,
	})
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
	if params.Keyword != "" {
		// 添加关键字搜索数量
		_, _ = helper.RedisGetClient().ZIncrBy(cs.ProductSearchHotKeywords, 1, params.Keyword).Result()
	}
	count := int64(-1)
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
		Products service.Products `json:"products"`
		Count    int64            `json:"count"`
	}{
		result,
		count,
	}
	return
}

// findByID find product by id
func (ctrl productCtrl) findByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	data, err := productSrv.FindByID(id)
	if err != nil {
		return
	}
	c.CacheMaxAge("1m")
	c.Body = data
	return
}

// updateByID update product by id
func (ctrl productCtrl) updateByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
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
		Specs:      params.Specs,
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
		Supplier:   params.Supplier,
		Rank:       params.Rank,
	}
	err = productSrv.UpdateByID(id, product)
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

	cat, err := productSrv.AddCategory(service.ProductCategory{
		Name:    params.Name,
		Level:   params.Level,
		Status:  params.Status,
		Belongs: params.Belongs,
		Rank:    params.Rank,
		Icon:    params.Icon,
	})
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
	count := int64(-1)
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
		ProductCategories service.ProductCategories `json:"productCategories,omitempty"`
		Count             int64                     `json:"count,omitempty"`
	}{
		result,
		count,
	}
	return
}

// updateCategoryByID update category by id
func (ctrl productCtrl) updateCategoryByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
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
		Rank:    params.Rank,
		Icon:    params.Icon,
	}
	err = productSrv.UpdateCategoryByID(id, cat)
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// findCategoryByID find category by id
func (ctrl productCtrl) findCategoryByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	data, err := productSrv.FindCategoryByID(id)
	if err != nil {
		return
	}
	c.CacheMaxAge("1m")
	c.Body = data
	return
}

// getMainImage get main image of product
func (ctrl productCtrl) getMainImage(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	quality, _ := strconv.Atoi(c.Param("quality"))
	width, _ := strconv.Atoi(c.Param("width"))
	height, _ := strconv.Atoi(c.Param("height"))
	product, err := productSrv.FindByID(id)
	if err != nil {
		return
	}
	file := product.Pics[0]
	if product.MainPic < len(product.Pics) {
		file = product.Pics[product.MainPic]
	}
	arr := strings.Split(file, "/")
	data, header, err := imageSrv.GetImageFromBucket(arr[len(arr)-2], arr[len(arr)-1], service.ImageOptimParams{
		Type:    c.Param("ext"),
		Quality: quality,
		Width:   width,
		Height:  height,
	})
	if err != nil {
		return
	}
	// 客户端缓存一周，缓存服务器缓存10分钟
	c.CacheMaxAge("168h", "10m")
	for k, values := range header {
		for _, v := range values {
			c.AddHeader(k, v)
		}
	}
	c.BodyBuffer = bytes.NewBuffer(data)
	return
}

// listSearchHotKeywords 热门搜索关键字
func (productCtrl) listSearchHotKeywords(c *elton.Context) (err error) {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 5
	}
	if limit > 10 {
		err = hes.New("热门关键字数量查询不能大于10")
		return
	}
	result, err := helper.RedisGetClient().ZRevRangeByScore(cs.ProductSearchHotKeywords, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   "+inf",
		Count: int64(limit),
	}).Result()
	if err != nil {
		return
	}
	c.CacheMaxAge("5m")
	c.Body = map[string][]string{
		"keywords": result,
	}
	return
}
