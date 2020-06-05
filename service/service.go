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

package service

import (
	"sort"

	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/helper"
	"github.com/vicanso/origin/log"
)

type (
	StatusInfo struct {
		Name  string `json:"name,omitempty"`
		Value int    `json:"value,omitempty"`
	}
	StatusInfoList []*StatusInfo
)

var (
	pgCreate    = helper.PGCreate
	pgGetClient = helper.PGGetClient
	pgQuery     = helper.PGQuery
	pgCount     = helper.PGCount

	logger = log.Default()

	redisSrv   = new(helper.Redis)
	productSrv = new(ProductSrv)
	regionSrv  = new(RegionSrv)
	brandSrv   = new(BrandSrv)

	statusInfoList StatusInfoList
	statusDict     map[int]string
)

func init() {
	statusDict = map[int]string{
		cs.StatusEnabled:  "启用",
		cs.StatusDisabled: "禁用",
	}
	statusInfoList = make(StatusInfoList, 0)
	for k, v := range statusDict {
		statusInfoList = append(statusInfoList, &StatusInfo{
			Name:  v,
			Value: k,
		})
	}
	sort.Slice(statusInfoList, func(i, j int) bool {
		return statusInfoList[i].Value < statusInfoList[j].Value
	})
}

func getStatusDesc(status int) string {
	value, ok := statusDict[status]
	if !ok {
		return ""
	}
	return value
}

func GetStatusList() []*StatusInfo {
	return statusInfoList
}
