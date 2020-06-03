// Copyright 2019 tree xie
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
	"strings"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/router"
	"github.com/vicanso/origin/service"
	"github.com/vicanso/origin/validate"
)

type (
	configurationCtrl      struct{}
	addConfigurationParams struct {
		Name      string     `json:"name,omitempty" validate:"xConfigName"`
		Category  string     `json:"category,omitempty" validate:"xConfigCategory"`
		Status    int        `json:"status,omitempty" validate:"xStatus"`
		Data      string     `json:"data,omitempty" validate:"xConfigData"`
		BeginDate *time.Time `json:"beginDate,omitempty"`
		EndDate   *time.Time `json:"endDate,omitempty"`
	}
	updateConfigurationParams struct {
		Status    int        `json:"status,omitempty" validate:"omitempty,xStatus"`
		Category  string     `json:"category,omitempty" validate:"omitempty,xConfigCategory"`
		Data      string     `json:"data,omitempty" validate:"omitempty,xConfigData"`
		BeginDate *time.Time `json:"beginDate,omitempty"`
		EndDate   *time.Time `json:"endDate,omitempty"`
	}
	listConfigurationParmas struct {
		listParams

		Name     string `json:"name,omitempty" validate:"omitempty,xConfigName"`
		Category string `json:"category,omitempty" validate:"omitempty,xConfigCategory"`
	}
)

func init() {
	g := router.NewGroup("/configurations", loadUserSession, shouldBeSu)
	ctrl := configurationCtrl{}

	g.GET(
		"/v1",
		ctrl.list,
	)

	g.POST(
		"/v1",
		newTracker(cs.ActionConfigurationAdd),
		ctrl.add,
	)
	g.GET(
		"/v1/{id}",
		ctrl.findByID,
	)
	g.PATCH(
		"/v1/{id}",
		newTracker(cs.ActionConfigurationUpdate),
		ctrl.updateByID,
	)
	g.DELETE(
		"/v1/{id}",
		newTracker(cs.ActionConfigurationDelete),
		ctrl.delete,
	)
}

func (params listConfigurationParmas) toConditions() []interface{} {
	conds := queryConditions{}
	if params.Name != "" {
		names := strings.Split(params.Name, ",")
		if len(names) > 1 {
			conds.add("name in (?)", names)
		} else {
			conds.add("name = (?)", names[0])
		}
	}

	if params.Category != "" {
		categories := strings.Split(params.Category, ",")
		if len(categories) > 1 {
			conds.add("category in (?)", categories)
		} else {
			conds.add("category = ?", categories[0])
		}
	}
	return conds.toArray()
}

// list configuration
func (ctrl configurationCtrl) list(c *elton.Context) (err error) {
	params := listConfigurationParmas{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	result, err := configSrv.List(params.toPGQueryParams(), params.toConditions()...)
	if err != nil {
		return
	}
	c.Body = map[string]interface{}{
		"configs": result,
	}
	return
}

// add configuration
func (ctrl configurationCtrl) add(c *elton.Context) (err error) {
	params := addConfigurationParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	us := getUserSession(c)
	conf, err := configSrv.Add(service.Configuration{
		Name:      params.Name,
		Category:  params.Category,
		Status:    params.Status,
		Data:      params.Data,
		Owner:     us.GetAccount(),
		BeginDate: params.BeginDate,
		EndDate:   params.EndDate,
	})
	if err != nil {
		return
	}
	c.Created(conf)
	return
}

// updateByID configuration
func (ctrl configurationCtrl) updateByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	params := updateConfigurationParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	err = configSrv.UpdateByID(id, service.Configuration{
		Status:    params.Status,
		Data:      params.Data,
		Category:  params.Category,
		BeginDate: params.BeginDate,
		EndDate:   params.EndDate,
	})
	if err != nil {
		return
	}

	c.NoContent()
	return
}

// delete configuration
func (ctrl configurationCtrl) delete(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	err = configSrv.DeleteByID(id)
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// findByID find configuration by id
func (ctrl configurationCtrl) findByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	data, err := configSrv.FindByID(id)
	if err != nil {
		return
	}
	c.Body = data
	return
}
