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
	"github.com/vicanso/origin/router"
	"github.com/vicanso/origin/service"
	"github.com/vicanso/origin/validate"
)

type (
	orderCommissionCtrl struct{}

	listOrderCommissionParams struct {
		listParams

		Recommender uint `json:"recommender,omitempty" validate:"omitempty,xUserID"`
	}

	listOrderCommissionsResp struct {
		Count            int64                    `json:"count,omitempty"`
		OrderCommissions service.OrderCommissions `json:"orderCommissions,omitempty"`
	}
)

func init() {
	ctrl := orderCommissionCtrl{}
	g := router.NewGroup("/order-commissions", loadUserSession, shouldBeLogined)

	g.GET(
		"/v1",
		ctrl.list,
	)
}

func (params listOrderCommissionParams) toConditions() []interface{} {
	conds := queryConditions{}
	if params.Recommender != 0 {
		conds.add("recommender = ?", params.Recommender)
	}
	return conds.toArray()
}

func (ctrl orderCommissionCtrl) list(c *elton.Context) (err error) {
	params := listOrderCommissionParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	us := getUserSession(c)
	params.Recommender = us.GetID()
	args := params.toConditions()
	queryParams := params.toPGQueryParams()
	count := int64(-1)
	if queryParams.Offset == 0 {
		count, err = orderCommissionSrv.Count(args...)
		if err != nil {
			return
		}
	}
	orderCommissions, err := orderCommissionSrv.List(queryParams, args...)
	if err != nil {
		return
	}
	c.Body = &listOrderCommissionsResp{
		Count:            count,
		OrderCommissions: orderCommissions,
	}
	return
}
