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
	"fmt"
	"net/http"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/router"
	"github.com/vicanso/origin/validate"
)

type (
	orderCtrl struct{}

	addOrderParams struct {
		Products []struct {
			ID    uint `json:"id,omitempty" validate:"xOrderProductID"`
			Count uint `json:"count,omitempty" validate:"xOrderProductCount"`
		} `json:"products,omitempty"`
	}
)

const (
	errOrderCtrlCategory = "order-ctrl"
)

var (
	errProductsIsEmpty = &hes.Error{
		Message:    "产品不能为空",
		StatusCode: http.StatusBadRequest,
		Category:   errOrderCtrlCategory,
	}
)

func init() {
	ctrl := orderCtrl{}
	g := router.NewGroup("/orders", loadUserSession, shouldLogined)

	g.POST(
		"/v1",
		newTracker(cs.ActionOrderAdd),
		ctrl.add,
	)
}

func (orderCtrl) add(c *elton.Context) (err error) {
	params := addOrderParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	if len(params.Products) == 0 {
		err = errProductsIsEmpty
		return
	}
	// us := getUserSession(c)
	// orderSrv.CreateWithSubOrders(us.GetAccount())

	fmt.Println(params)
	return
}
