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
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/origin/config"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/helper"
	"github.com/vicanso/origin/util"
	"gorm.io/gorm"
)

const (
	mockTimeKey               = "mockTime"
	sessionSignedKeyCateogry  = "signedKey"
	blockIPCategory           = "blockIP"
	routerConfigCategory      = "router"
	routerConcurrencyCategory = "routerConcurrency"
	// 佣金配置
	orderCommissionCategory = "orderCommission"
	// 营销分组
	marketingGroupCategory = "marketingGroup"
)

var (
	signedKeys = new(elton.RWMutexSignedKeys)

	defaultMarketingGroups = new(MarketingGroups)
)

type (
	Configurations []*Configuration
	// Configuration configuration of application
	Configuration struct {
		helper.Model

		// 配置名称，唯一
		Name string `json:"name,omitempty" gorm:"type:varchar(30);not null;unique_index:idx_configuration_name"`
		// 配置分类
		Category string `json:"category,omitempty" gorm:"type:varchar(20)"`
		// 配置由谁创建
		Owner string `json:"owner,omitempty" gorm:"type:varchar(20);not null"`
		// 配置状态
		Status     int    `json:"status,omitempty" gorm:"index:idx_configuration_status"`
		StatusDesc string `json:"statusDesc,omitempty" gorm:"-"`
		Data       string `json:"data,omitempty"`
		// 启用开始时间
		BeginDate *time.Time `json:"beginDate,omitempty"`
		// 启用结束时间
		EndDate *time.Time `json:"endDate,omitempty"`
	}
	// ConfigurationSrv configuration service
	ConfigurationSrv struct {
	}

	MarketingGroup struct {
		Name  string `json:"name,omitempty"`
		Owner uint   `json:"owner,omitempty"`
	}
	// MarketingGroups marketing groups
	MarketingGroups struct {
		sync.RWMutex
		groups []*MarketingGroup
	}
)

func init() {
	err := helper.PGAutoMigrate(&Configuration{})
	if err != nil {
		panic(err)
	}
	signedKeys.SetKeys(config.GetSignedKeys())
}

func (cs Configurations) AfterFind(tx *gorm.DB) (err error) {
	for _, c := range cs {
		err = c.AfterFind(tx)
		if err != nil {
			return
		}
	}
	return
}

func (c *Configuration) AfterFind(_ *gorm.DB) (err error) {
	c.StatusDesc = getStatusDesc(c.Status)
	return
}

// IsValid check the config is valid
func (conf *Configuration) IsValid() bool {
	if conf.Status != cs.StatusEnabled {
		return false
	}
	return util.IsBetween(conf.BeginDate, conf.EndDate)
}

// Set set groups
func (mg *MarketingGroups) Set(data []string) {
	groups := make([]*MarketingGroup, 0)
	for _, item := range data {
		group := &MarketingGroup{}
		_ = json.Unmarshal([]byte(item), group)
		if group.Name != "" {
			groups = append(groups, group)
		}
	}
	mg.Lock()
	defer mg.Unlock()
	mg.groups = groups
}

// List list marketing groups
func (mg *MarketingGroups) List() (groups []*MarketingGroup) {
	mg.RLock()
	defer mg.RUnlock()
	if len(mg.groups) == 0 {
		return nil
	}
	list := mg.groups[0:]
	return list
}

// GetOwner get group owner
func (mg *MarketingGroups) GetOwner(group string) (owner uint) {
	mg.RLock()
	defer mg.RUnlock()
	for _, item := range mg.groups {
		if item.Name == group {
			owner = item.Owner
			break
		}
	}
	return
}

// ListMarketingGroup list marketing group
func ListMarketingGroup() (groups []*MarketingGroup) {
	return defaultMarketingGroups.List()
}

// createByID create a configuration by id
func (srv *ConfigurationSrv) createByID(id uint) *Configuration {
	c := &Configuration{}
	c.Model.ID = id
	return c
}

// Add add configuration
func (srv *ConfigurationSrv) Add(data Configuration) (conf *Configuration, err error) {
	conf = &data
	err = pgCreate(conf)
	return
}

// UpdateByID update configuration by id
func (srv *ConfigurationSrv) UpdateByID(id uint, conf Configuration) (err error) {
	err = pgGetClient().Model(srv.createByID(id)).Updates(conf).Error
	return
}

// FindByID find configuration by id
func (srv *ConfigurationSrv) FindByID(id uint) (config *Configuration, err error) {
	config = new(Configuration)
	err = pgGetClient().First(config, "id = ?", id).Error
	return
}

// Available get available configs
func (srv *ConfigurationSrv) Available() (configs []*Configuration, err error) {
	configs = make([]*Configuration, 0)
	now := time.Now()
	err = pgQuery(PGQueryParams{
		Order: "-updatedAt",
	}, "status = ? and begin_date < ? and end_date > ?", cs.StatusEnabled, now, now).Find(&configs).Error
	if err != nil {
		return
	}
	return
}

// Refresh refresh configurations
func (srv *ConfigurationSrv) Refresh() (err error) {
	configs, err := srv.Available()
	if err != nil {
		return
	}
	var mockTimeConfig *Configuration

	routerConfigs := make([]*Configuration, 0)
	var signedKeysConfig *Configuration
	blockIPList := make([]string, 0)
	routerConcurrencyConfigs := make([]string, 0)
	orderCommissionConfigs := make([]string, 0)
	groupConfigs := make([]string, 0)

	for _, item := range configs {
		if item.Name == mockTimeKey {
			mockTimeConfig = item
			continue
		}

		switch item.Category {
		case routerConfigCategory:
			// 路由配置
			routerConfigs = append(routerConfigs, item)
		case sessionSignedKeyCateogry:
			// session的签名串配置
			signedKeysConfig = item
		case blockIPCategory:
			blockIPList = append(blockIPList, item.Data)
		case routerConcurrencyCategory:
			routerConcurrencyConfigs = append(routerConcurrencyConfigs, item.Data)
		case orderCommissionCategory:
			orderCommissionConfigs = append(orderCommissionConfigs, item.Data)
		case marketingGroupCategory:
			groupConfigs = append(groupConfigs, item.Data)
		}
	}

	// 如果未配置mock time，则设置为空
	if mockTimeConfig == nil {
		util.SetMockTime("")
	} else {
		util.SetMockTime(mockTimeConfig.Data)
	}

	// 如果数据库中未配置，则使用默认配置
	if signedKeysConfig == nil {
		signedKeys.SetKeys(config.GetSignedKeys())
	} else {
		keys := strings.Split(signedKeysConfig.Data, ",")
		signedKeys.SetKeys(keys)
	}

	defaultOrderCommissionConfigs.Set(orderCommissionConfigs)

	defaultMarketingGroups.Set(groupConfigs)

	// 更新router configs
	updateRouterConfigs(routerConfigs)

	ResetIPBlocker(blockIPList)
	ResetRouterConcurrency(routerConcurrencyConfigs)
	return
}

// GetSignedKeys get signed keys
func GetSignedKeys() elton.SignedKeysGenerator {
	return signedKeys
}

// List list configurations
func (srv *ConfigurationSrv) List(params PGQueryParams, args ...interface{}) (result Configurations, err error) {
	result = make(Configurations, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// DeleteByID delete configuration
func (srv *ConfigurationSrv) DeleteByID(id uint) (err error) {
	err = pgGetClient().Unscoped().Delete(srv.createByID(id)).Error
	return
}
