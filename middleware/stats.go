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

package middleware

import (
	"github.com/dustin/go-humanize"
	"github.com/vicanso/elton"
	M "github.com/vicanso/elton/middleware"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/helper"
	"github.com/vicanso/origin/log"
	"github.com/vicanso/origin/util"
	"go.uber.org/zap"
)

func NewStats() elton.Handler {
	logger := log.Default()
	return M.NewStats(M.StatsConfig{
		OnStats: func(info *M.StatsInfo, c *elton.Context) {
			// ping 的日志忽略
			if info.URI == "/ping" {
				return
			}
			sid := util.GetSessionID(c)
			logger.Info("access log",
				zap.String("uuid", c.GetRequestHeader("X-UUID")),
				zap.String("id", info.CID),
				zap.String("ip", info.IP),
				zap.String("sid", sid),
				zap.String("method", info.Method),
				zap.String("route", info.Route),
				zap.String("uri", info.URI),
				zap.Int("status", info.Status),
				zap.Uint32("connecting", info.Connecting),
				zap.String("consuming", info.Consuming.String()),
				zap.String("size", humanize.Bytes(uint64(info.Size))),
				zap.Int("bytes", info.Size),
			)
			tags := map[string]string{
				"method": info.Method,
				"route":  info.Route,
			}
			fields := map[string]interface{}{
				"id":         info.CID,
				"ip":         info.IP,
				"sid":        sid,
				"uri":        info.URI,
				"status":     info.Status,
				"use":        info.Consuming.Milliseconds(),
				"size":       info.Size,
				"connecting": info.Connecting,
			}
			helper.GetInfluxSrv().Write(cs.MeasurementHTTP, fields, tags)
		},
	})
}
