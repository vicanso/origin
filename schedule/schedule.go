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

package schedule

import (
	"time"

	"github.com/vicanso/origin/cs"

	"github.com/robfig/cron/v3"
	"github.com/vicanso/origin/helper"
	"github.com/vicanso/origin/log"
	"github.com/vicanso/origin/service"

	"go.uber.org/zap"
)

func init() {
	c := cron.New()
	_, _ = c.AddFunc("@every 5m", redisCheck)
	_, _ = c.AddFunc("@every 1m", configRefresh)
	_, _ = c.AddFunc("@every 5m", redisStats)
	_, _ = c.AddFunc("@every 1m", pgStats)
	_, _ = c.AddFunc("00 00 * * *", resetProductSearchHotKeywords)
	// 测试暂时每5分钟自动生成
	_, _ = c.AddFunc("@every 5m", generateOrderCommission)
	go func() {
		time.Sleep(time.Second)
		generateOrderCommission()
	}()
	c.Start()
}

func redisCheck() {
	err := helper.RedisPing()
	if err != nil {
		log.Default().Error("redis check fail",
			zap.Error(err),
		)
		service.AlarmError("redis check fail")
	}
}

func configRefresh() {
	configSrv := new(service.ConfigurationSrv)
	err := configSrv.Refresh()
	if err != nil {
		log.Default().Error("config refresh fail",
			zap.Error(err),
		)
		service.AlarmError("config refresh fail, " + err.Error())
	}
}

func redisStats() {
	stats := helper.RedisStats()
	helper.GetInfluxSrv().Write(cs.MeasurementRedisStats, stats, nil)
}

func pgStats() {
	stats := helper.PGStats()
	helper.GetInfluxSrv().Write(cs.MeasurementPGStats, stats, nil)
}

func resetProductSearchHotKeywords() {
	_, err := helper.RedisGetClient().ZRemRangeByRank(cs.ProductSearchHotKeywords, 0, -1).Result()
	if err != nil {
		log.Default().Error("reset product search hot key words fail",
			zap.Error(err),
		)
	}
}

func generateOrderCommission() {
	orderCommissionSrv := new(service.OrderCommissionSrv)
	err := orderCommissionSrv.Do()
	if err != nil {
		log.Default().Error("order commission generate fail",
			zap.Error(err),
		)
		service.AlarmError("order commission geerate fail, " + err.Error())
	}
}
