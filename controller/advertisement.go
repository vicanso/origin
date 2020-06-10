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
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/router"
	"github.com/vicanso/origin/service"
	"github.com/vicanso/origin/validate"
)

type (
	addAdvertisementParams struct {
		Status    int        `json:"status,omitempty" validate:"xStatus"`
		Link      string     `json:"link,omitempty" validate:"xAdvertisementLink"`
		Summary   string     `json:"summary,omitempty" validate:"xAdvertisementSummary"`
		Category  string     `json:"category,omitempty" validate:"xAdvertisementCategory"`
		StartedAt *time.Time `json:"startedAt,omitempty" validate:"required"`
		EndedAt   *time.Time `json:"endedAt,omitempty" validate:"required"`
	}
	updateAdvertisementParams struct {
		Status    int        `json:"status,omitempty" validate:"omitempty,xStatus"`
		Link      string     `json:"link,omitempty" validate:"omitempty,xAdvertisementLink"`
		Summary   string     `json:"summary,omitempty" validate:"omitempty,xAdvertisementSummary"`
		Category  string     `json:"category,omitempty" validate:"omitempty,xAdvertisementCategory"`
		StartedAt *time.Time `json:"startedAt,omitempty"`
		EndedAt   *time.Time `json:"endedAt,omitempty"`
	}
	listAdvertisementParams struct {
		listParams

		Status string `json:"status,omitempty" validate:"omitempty,xStatus"`
	}

	advertisementCtrl struct{}
)

func init() {
	g := router.NewGroup("/advertisements")
	ctrl := advertisementCtrl{}

	g.GET(
		"/v1",
		ctrl.list,
	)
	g.POST(
		"/v1",
		loadUserSession,
		newTracker(cs.ActionAdvertisementAdd),
		checkMarketingGroup,
		ctrl.add,
	)
	g.PATCH(
		"/v1/{id}",
		loadUserSession,
		newTracker(cs.ActionAdvertisementUpdate),
		checkMarketingGroup,
		ctrl.updateByID,
	)
}

func (params listAdvertisementParams) toConditions() (conditions []interface{}) {
	conds := queryConditions{}
	if params.Status != "" {
		conds.add("status = ?", params.Status)
	}
	return conds.toArray()
}

// add add advertisement
func (ctrl advertisementCtrl) add(c *elton.Context) (err error) {
	params := addAdvertisementParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}

	advertisement, err := advertisementSrv.Add(service.Advertisement{
		Status:    params.Status,
		Link:      params.Link,
		Summary:   params.Summary,
		Category:  params.Category,
		StartedAt: params.StartedAt,
		EndedAt:   params.EndedAt,
	})
	if err != nil {
		return
	}
	c.Created(advertisement)
	return
}

// updateByID update advertisement by id
func (ctrl advertisementCtrl) updateByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	params := updateAdvertisementParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	err = advertisementSrv.UpdateByID(id, service.Advertisement{
		Status:    params.Status,
		Link:      params.Link,
		Summary:   params.Summary,
		Category:  params.Category,
		StartedAt: params.StartedAt,
		EndedAt:   params.EndedAt,
	})
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// findByID find advertisement by id
func (ctrl advertisementCtrl) findByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	data, err := advertisementSrv.FindByID(id)
	if err != nil {
		return
	}
	c.CacheMaxAge("1m")
	c.Body = data
	return
}

func (ctrl advertisementCtrl) list(c *elton.Context) (err error) {
	params := listAdvertisementParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	count := -1
	args := params.toConditions()
	queryParams := params.toPGQueryParams()
	if queryParams.Offset == 0 {
		count, err = advertisementSrv.Count(args...)
		if err != nil {
			return
		}
	}
	result, err := advertisementSrv.List(queryParams, args...)
	if err != nil {
		return
	}
	c.CacheMaxAge("1m")
	c.Body = &struct {
		Advertisements []*service.Advertisement `json:"advertisements,omitempty"`
		Count          int                      `json:"count,omitempty"`
	}{
		result,
		count,
	}
	return
}
