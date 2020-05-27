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

	listProvince struct {
		listParams
	}
)

func init() {
	ctrl := regionCtrl{}
	g := router.NewGroup("/regions")

	g.POST(
		"/v1/import/{category}",
		loadUserSession,
		shouldBeSu,
		ctrl.importFromFile,
	)

	g.GET(
		"/v1/provinces",
		ctrl.listProvince,
	)
	g.GET(
		"/v1/cities/{province}",
		ctrl.listCity,
	)
}

func (ctrl *regionCtrl) importFromFile(c *elton.Context) (err error) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		return
	}
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	arr := make([]map[string]string, 0)
	err = json.Unmarshal(buf, &arr)
	if err != nil {
		return
	}

	for _, item := range arr {
		code, e := strconv.Atoi(item["code"])
		if e != nil {
			err = e
			return
		}
		name := item["name"]
		switch c.Param("category") {
		case cs.RegionProvince:
			err = regionSrv.AddProvince(&service.RegionProvince{
				Name:    name,
				Code:    code,
				Country: service.China,
			})
			if err != nil {
				return
			}
		case cs.RegionCity:

			province, e := strconv.Atoi(item["provinceCode"])
			if e != nil {
				err = e
				return
			}
			err = regionSrv.AddCity(&service.RegionCity{
				Name:     name,
				Code:     code,
				Province: province,
			})
			if err != nil {
				return
			}

		case cs.RegionArea:
			city, e := strconv.Atoi(item["cityCode"])
			if e != nil {
				err = e
				return
			}
			err = regionSrv.AddArea(&service.RegionArea{
				Name: name,
				Code: code,
				City: city,
			})
			if err != nil {
				return
			}
		case cs.RegionStreet:
			area, e := strconv.Atoi(item["areaCode"])
			if e != nil {
				err = e
				return
			}
			err = regionSrv.AddStreet(&service.RegionStreet{
				Name: name,
				Code: code,
				Area: area,
			})
			if err != nil {
				return
			}

		default:
			err = hes.New("category is invalid")
		}
	}
	if err != nil {
		return
	}
	c.NoContent()
	return
}

func (ctrl *regionCtrl) listProvince(c *elton.Context) (err error) {
	params := listProductParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	data, err := regionSrv.ListProvince(params.toPGQueryParams())
	if err != nil {
		return
	}
	c.CacheMaxAge("5m")
	c.Body = map[string][]*service.RegionProvince{
		"provinces": data,
	}
	return
}

func (ctrl *regionCtrl) listCity(c *elton.Context) (err error) {
	params := listProductParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	province := c.Param("province")
	data, err := regionSrv.ListCity(params.toPGQueryParams(), "province = ?", province)
	if err != nil {
		return
	}
	c.CacheMaxAge("5m")
	c.Body = map[string][]*service.RegionCity{
		"cities": data,
	}
	return
}
