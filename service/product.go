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
	"net/http"
	"time"

	"github.com/lib/pq"
	"github.com/vicanso/hes"
	lruTTL "github.com/vicanso/lru-ttl"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/helper"
	"github.com/vicanso/origin/util"
)

type (
	Product struct {
		helper.Model

		Name string `json:"name,omitempty" gorm:"type:varchar(30);not null;index:idx_product_name"`
		// 单价
		Price float64 `json:"price,omitempty" gorm:"not null"`
		// 规格，规格+单位表示完整的购买单位，如规则为250，单位为克
		Specs uint `json:"specs,omitempty" gorm:"not null"`
		// 单位
		Unit    string         `json:"unit,omitempty" gorm:"not null"`
		Catalog string         `json:"catalog,omitempty"`
		Pics    pq.StringArray `json:"pics,omitempty" gorm:"type:text[]"`
		// 主图，从1开始
		MainPic int    `json:"mainPic,omitempty"`
		SN      string `json:"sn,omitempty"`

		Status     int    `json:"status,omitempty" gorm:"index:idx_product_status"`
		StatusDesc string `json:"statusDesc,omitempty" gorm:"-"`

		// 排序
		Rank int `json:"rank,omitempty"`
		// 关键字
		Keywords string `json:"keywords,omitempty"`

		Categories pq.Int64Array `json:"categories,omitempty" gorm:"type:int[]"`
		// 产品分类说明（在获取数据后转换生成）
		CategoriesDesc []string   `json:"categoriesDesc,omitempty" gorm:"-"`
		StartedAt      *time.Time `json:"startedAt,omitempty" gorm:"not null;index:idx_product_started_at"`
		EndedAt        *time.Time `json:"endedAt,omitempty" gorm:"not null;index:idx_product_ended_at"`

		// 产地
		Origin string `json:"origin,omitempty"`
		// 产地说明
		OrginDesc string `json:"orginDesc,omitempty" gorm:"-"`

		// 产品品牌
		Brand uint `json:"brand,omitempty"`
		// 产品品牌说明
		BrandDesc string `json:"brandDesc,omitempty" gorm:"-"`

		// 供应商
		Supplier uint `json:"supplier,omitempty"`
		// 供应商说明
		SupplierDesc string `json:"supplierDesc,omitempty" gorm:"-"`

		// 是否有效(是否可购买)
		Available bool `json:"available,omitempty" gorm:"-"`
	}
	// ProductCategory product category
	ProductCategory struct {
		helper.Model

		Name       string `json:"name,omitempty" grom:"not null;unique_index:idx_product_category_name"`
		Level      int    `json:"level,omitempty" grom:"not null;index:idx_product_category_level"`
		Status     int    `json:"status,omitempty" gorm:"not null;index:idx_product_category_status"`
		StatusDesc string `json:"statusDesc,omitempty" gorm:"-"`
		// 所属分类
		Belongs pq.Int64Array `json:"belongs,omitempty" gorm:"type:int[]"`
		// 所属分类描述
		BelongsDesc []string `json:"belongsDesc,omitempty" gorm:"-"`
		// 排序
		Rank int `json:"rank,omitempty"`
		// 图标
		Icon string `json:"icon,omitempty"`
	}
	ProductSrv struct{}
)

const (
	errProductCategory = "product"
)

var (
	errProductUnavailable = &hes.Error{
		Message:    "产品状态非可销售状态或过已销售期",
		StatusCode: http.StatusBadRequest,
		Category:   errProductCategory,
	}
)

var (
	// productCategoryNameCache product category's name cache
	productCategoryNameCache *lruTTL.Cache
)

func init() {
	ttl := 300 * time.Second
	// 本地开发环境，设置缓存为1秒
	if util.IsDevelopment() {
		ttl = time.Second
	}
	productCategoryNameCache = lruTTL.New(1000, ttl)

	pgGetClient().AutoMigrate(
		&Product{},
		&ProductCategory{},
	)
}

func (p *Product) AfterFind() (err error) {
	p.StatusDesc = getStatusDesc(p.Status)

	// 根据产品分类转换对应分类名称
	categoriesDesc := make([]string, 0)
	for _, id := range p.Categories {
		// 如果获取失败，忽略出错
		name, _ := productSrv.GetCategoryNameFromCache(uint(id))
		if name != "" {
			categoriesDesc = append(categoriesDesc, name)
		}
	}
	p.CategoriesDesc = categoriesDesc

	// 如果获取失败，忽略出错
	p.OrginDesc, _ = regionSrv.GetNameFromCache(p.Origin, 0)

	// 如果获取失败，忽略出错
	p.BrandDesc, _ = brandSrv.GetNameFromCache(p.Brand)
	p.Available = p.IsAvailable()

	return
}

func (pc *ProductCategory) AfterFind() (err error) {
	pc.StatusDesc = getStatusDesc(pc.Status)

	// 生成上级分类描述
	belongsDesc := make([]string, 0)
	for _, id := range pc.Belongs {
		// 如果获取失败，忽略出错
		name, _ := productSrv.GetCategoryNameFromCache(uint(id))
		if name != "" {
			belongsDesc = append(belongsDesc, name)
		}
	}
	pc.BelongsDesc = belongsDesc

	return
}

func (pc *ProductCategory) BeforeCreate() (err error) {
	// 排序默认设置为1
	if pc.Rank == 0 {
		pc.Rank = 1
	}
	return
}

func (srv *ProductSrv) createByID(id uint) *Product {
	p := &Product{}
	p.Model.ID = id
	return p
}

// Available 是否可购买
func (product *Product) IsAvailable() bool {
	if product.Status != cs.StatusEnabled {
		return false
	}
	return util.IsBetween(product.StartedAt, product.EndedAt)
}

// CheckAvailable check product is available
func (product *Product) CheckAvailable() error {
	if !product.IsAvailable() {
		return errProductUnavailable
	}
	return nil
}

// Add add product
func (srv *ProductSrv) Add(data Product) (product *Product, err error) {
	product = &data
	err = pgCreate(product)
	return
}

// UpdateByID update product by id
func (srv *ProductSrv) UpdateByID(id uint, product Product) (err error) {
	err = pgGetClient().Model(srv.createByID(id)).Updates(product).Error
	return
}

// FindByID find product by id
func (srv *ProductSrv) FindByID(id uint) (product *Product, err error) {
	product = new(Product)
	err = pgGetClient().First(product, "id = ?", id).Error
	return
}

// List list product
func (srv *ProductSrv) List(params PGQueryParams, args ...interface{}) (result []*Product, err error) {
	result = make([]*Product, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// Count count the product
func (srv *ProductSrv) Count(args ...interface{}) (count int, err error) {
	return pgCount(&Product{}, args...)
}

func (srv *ProductSrv) createCategoryByID(id uint) *ProductCategory {
	c := &ProductCategory{}
	c.Model.ID = id
	return c
}

// AddCategory add category
func (srv *ProductSrv) AddCategory(data ProductCategory) (cat *ProductCategory, err error) {
	cat = &data
	err = pgCreate(cat)
	return
}

// UpdateCategoryByID update category by id
func (srv *ProductSrv) UpdateCategoryByID(id uint, attrs ...interface{}) (err error) {
	err = pgGetClient().Model(srv.createCategoryByID(id)).Update(attrs...).Error
	return
}

// FindCategoryByID find category by id
func (srv *ProductSrv) FindCategoryByID(id uint) (cat *ProductCategory, err error) {
	cat = new(ProductCategory)
	err = pgGetClient().First(cat, "id = ?", id).Error
	return
}

// ListCategory list category
func (srv *ProductSrv) ListCategory(params PGQueryParams, args ...interface{}) (result []*ProductCategory, err error) {
	result = make([]*ProductCategory, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// CountCategory count category
func (srv *ProductSrv) CountCategory(args ...interface{}) (count int, err error) {
	return pgCount(&ProductCategory{}, args...)
}

// GetCategoryNameFromCache get product category's name from cache
// If not exists, it will get from db and save to cache
func (srv *ProductSrv) GetCategoryNameFromCache(id uint) (name string, err error) {
	if id == 0 {
		return
	}
	value, ok := productCategoryNameCache.Get(id)
	if ok {
		return value.(string), nil
	}
	cat, err := srv.FindCategoryByID(id)
	if err != nil {
		return
	}
	name = cat.Name
	productCategoryNameCache.Add(id, name)
	return
}
