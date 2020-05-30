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
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/router"
	"github.com/vicanso/origin/service"
	"github.com/vicanso/origin/validate"
)

type (
	regionCtrl struct{}

	listRegionParams struct {
		listParams

		Category string `json:"category,omitempty" validate:"omitempty,xRegionCategory"`
		Parent   string `json:"parent,omitempty" validate:"omitempty,xRegionParent"`
		Keyword  string `json:"keyword,omitempty" validate:"omitempty,xKeyword"`
		Status   string `json:"status,omitempty" validate:"omitempty,xStatus"`
	}

	updateRegionParams struct {
		Name   string `json:"name,omitempty" validate:"omitempty,xRegionName"`
		Status int    `json:"status,omitempty" validate:"omitempty,xStatus"`
	}
)

func init() {
	ctrl := regionCtrl{}
	g := router.NewGroup("/regions")

	g.GET(
		"/v1/categories",
		noCacheIfSetNoCache,
		ctrl.listCategory,
	)

	g.POST(
		"/v1/import/{category}",
		loadUserSession,
		shouldBeSu,
		ctrl.importFromFile,
	)

	g.GET(
		"/v1",
		ctrl.listRegion,
	)
	g.GET(
		"/v1/{id}",
		ctrl.findByID,
	)
	g.PATCH(
		"/v1/{id}",
		loadUserSession,
		checkMarketingGroup,
		ctrl.updateByID,
	)
}

func (params listRegionParams) toConditions() []interface{} {
	conds := queryConditions{}
	if params.Keyword != "" {
		conds.add("name ILIKE ?", "%"+params.Keyword+"%")
	}
	if params.Status != "" {
		conds.add("status = ?", params.Status)
	}
	if params.Category != "" {
		conds.add("category = ?", regionSrv.GetCategoryIndex(params.Category))
	}
	if params.Parent != "" {
		conds.add("parent = ?", params.Parent)
	}
	return conds.toArray()
}

func (ctrl regionCtrl) listCategory(c *elton.Context) (err error) {
	c.CacheMaxAge("5m")
	c.Body = map[string][]*service.RegionCategory{
		"categories": regionSrv.ListCategory(),
	}
	return
}

func (ctrl regionCtrl) importFromFile(c *elton.Context) (err error) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		return
	}
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	category := c.Param("category")
	categoryIndex := regionSrv.GetCategoryIndex(category)

	if cs.RegionCountry == category {
		m := make(map[string]string)
		err = json.Unmarshal(buf, &m)
		if err != nil {
			return
		}
		for key, value := range m {
			err = regionSrv.Add(&service.Region{
				Category: categoryIndex,
				Name:     value,
				Code:     key,
				Status:   cs.StatusDisabled,
			})

			if err != nil {
				return
			}
		}
		c.NoContent()
		return
	}

	arr := make([]map[string]string, 0)
	err = json.Unmarshal(buf, &arr)
	if err != nil {
		return
	}

	for _, item := range arr {
		region := service.Region{
			Category: categoryIndex,
			Name:     item["name"],
			Code:     item["code"],
			Status:   cs.StatusEnabled,
		}
		switch category {
		case cs.RegionProvince:
			region.Parent = "CN"
		case cs.RegionCity:
			region.Parent = item["provinceCode"]
		case cs.RegionArea:
			region.Parent = item["cityCode"]
		case cs.RegionStreet:
			region.Parent = item["areaCode"]
		default:
			err = hes.New("category is invalid")
			return
		}
		err = regionSrv.Add(&region)
		if err != nil {
			return
		}
	}
	c.NoContent()
	return
}

func (ctrl regionCtrl) listRegion(c *elton.Context) (err error) {
	params := listRegionParams{}
	query := c.Query()
	err = validate.Do(&params, query)
	if err != nil {
		return
	}
	queryParmas := params.toPGQueryParams()
	args := params.toConditions()
	count := -1
	if queryParmas.Offset == 0 {
		count, err = regionSrv.Count(args...)
		if err != nil {
			return
		}
	}
	result, err := regionSrv.List(queryParmas, args...)
	if err != nil {
		return
	}
	c.Body = map[string]interface{}{
		"count":   count,
		"regions": result,
	}
	return
}

func (ctrl regionCtrl) findByID(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	data, err := regionSrv.FindByID(uint(id))
	if err != nil {
		return
	}
	c.CacheMaxAge("1m")
	c.Body = data
	return
}

func (ctrl regionCtrl) updateByID(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	params := updateRegionParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	region := service.Region{
		Name:   params.Name,
		Status: params.Status,
	}
	err = regionSrv.UpdateByID(uint(id), region)
	if err != nil {
		return
	}
	c.NoContent()
	return
}
