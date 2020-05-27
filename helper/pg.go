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

package helper

import (
	"net/http"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/vicanso/hes"
	"github.com/vicanso/origin/config"
	"github.com/vicanso/origin/cs"
	"github.com/vicanso/origin/log"
	"go.uber.org/zap"
)

var (
	pgClient    *gorm.DB
	pgStatsHook *pgStats

	ErrPGTooManyQueryProcessing = &hes.Error{
		Message:    "too many query processing",
		StatusCode: http.StatusInternalServerError,
		Category:   "pg",
	}
	ErrPGTooManyUpdateProcessing = &hes.Error{
		Message:    "too many update processing",
		StatusCode: http.StatusInternalServerError,
		Category:   "pg",
	}
)

const (
	queryCMD  = "query"
	updateCMD = "update"
)

type (
	pgStats struct {
		slow                time.Duration
		maxQueryProcessing  uint32
		maxUpdateProcessing uint32
		queryProcessing     uint32
		updateProcessing    uint32
		total               uint64
	}

	// PGQueryParams pg query params
	PGQueryParams struct {
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
		Fields string `json:"fields"`
		Order  string `json:"order"`
	}

	Model struct {
		ID        uint       `gorm:"primary_key" json:"id"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt time.Time  `json:"updatedAt"`
		DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	}
)

func (ps *pgStats) getProcessingAndTotal() (uint32, uint32, uint64) {
	queryProcessing := atomic.LoadUint32(&ps.queryProcessing)
	updateProcessing := atomic.LoadUint32(&ps.updateProcessing)
	total := atomic.LoadUint64(&ps.total)
	return queryProcessing, updateProcessing, total
}

// Before before pg sql handle
func (ps *pgStats) Before(category string) (callback func(scope *gorm.Scope)) {
	return func(scope *gorm.Scope) {
		atomic.AddUint64(&ps.total, 1)

		switch category {
		case queryCMD:
			v := atomic.AddUint32(&ps.queryProcessing, 1)
			if v > ps.maxQueryProcessing {
				_ = scope.Err(ErrPGTooManyQueryProcessing)
			}
		case updateCMD:
			v := atomic.AddUint32(&ps.updateProcessing, 1)
			if v > ps.maxUpdateProcessing {
				_ = scope.Err(ErrPGTooManyUpdateProcessing)
			}
		}
		scope.InstanceSet(string(startedAtKey), time.Now())
	}
}

// After after pg sql handle
func (ps *pgStats) After(category string) func(*gorm.Scope) {
	return func(scope *gorm.Scope) {
		switch category {
		case queryCMD:
			atomic.AddUint32(&ps.queryProcessing, ^uint32(0))
		case updateCMD:
			atomic.AddUint32(&ps.updateProcessing, ^uint32(0))
		}

		value, ok := scope.InstanceGet(string(startedAtKey))
		if !ok {
			return
		}
		startedAt, ok := value.(time.Time)
		if !ok {
			return
		}
		use := time.Since(startedAt)
		db := scope.DB()
		if time.Since(startedAt) > ps.slow || db.Error != nil {
			message := ""
			if db.Error != nil {
				message = db.Error.Error()
			}
			logger.Info("pg process slow or error",
				zap.String("table", scope.TableName()),
				zap.String("category", category),
				zap.String("use", use.String()),
				zap.Int64("rowsAffected", db.RowsAffected),
				zap.String("error", message),
			)
			tags := map[string]string{
				"table":    scope.TableName(),
				"category": category,
			}
			fields := map[string]interface{}{
				"use":          use.Milliseconds(),
				"rowsAffected": db.RowsAffected,
				"error":        message,
			}
			GetInfluxSrv().Write(cs.MeasurementPG, fields, tags)
		}
	}
}

func init() {
	str := config.GetPostgresConnectString()
	pgConfig := config.GetPostgresConfig()
	reg := regexp.MustCompile(`password=\S*`)
	maskStr := reg.ReplaceAllString(str, "password=***")
	logger.Info("connect to pg",
		zap.String("args", maskStr),
	)
	db, err := gorm.Open("postgres", str)
	if err != nil {
		panic(err)
	}
	pgStatsHook = &pgStats{
		slow:                pgConfig.Slow,
		maxQueryProcessing:  pgConfig.MaxQueryProcessing,
		maxUpdateProcessing: pgConfig.MaxUpdateProcessing,
	}

	db.SetLogger(log.PGLogger())
	db.Callback().Query().Before("gorm:query").Register("stats:beforeQuery", pgStatsHook.Before(queryCMD))
	db.Callback().Query().After("gorm:query").Register("stats:afterQuery", pgStatsHook.After(queryCMD))
	db.Callback().Update().Before("gorm:update").Register("stats:beforeUpdate", pgStatsHook.Before(updateCMD))
	db.Callback().Update().After("gorm:update").Register("stats:afterUpdate", pgStatsHook.After(updateCMD))

	pgClient = db
}

// PGCreate pg create
func PGCreate(data interface{}) (err error) {
	err = pgClient.Create(data).Error
	return
}

// PGFirstOrCreate pg first of create
func PGFirstOrCreate(out interface{}, where ...interface{}) (err error) {
	err = pgClient.FirstOrCreate(out, where...).Error
	return err
}

// PGGetClient pg client
func PGGetClient() *gorm.DB {
	return pgClient
}

// PGFormatOrder format order
func PGFormatOrder(sort string) string {
	arr := strings.Split(sort, ",")
	newSort := []string{}
	for _, item := range arr {
		if item[0] == '-' {
			newSort = append(newSort, strcase.ToSnake(item[1:])+" desc")
		} else {
			newSort = append(newSort, strcase.ToSnake(item))
		}
	}
	return strings.Join(newSort, ",")
}

// PGFormatSelect format select
func PGFormatSelect(fields string) string {
	return strcase.ToSnake(fields)
}

// PGStats get pg stats
func PGStats() map[string]interface{} {
	queryProcessing, updateProcessing, total := pgStatsHook.getProcessingAndTotal()
	return map[string]interface{}{
		"queryProcessing":  queryProcessing,
		"updateProcessing": updateProcessing,
		"total":            total,
	}
}

// PGQuery pg query
func PGQuery(params PGQueryParams, args ...interface{}) *gorm.DB {
	db := PGGetClient()
	if params.Limit != 0 {
		db = db.Limit(params.Limit)
	}
	if params.Offset != 0 {
		db = db.Offset(params.Offset)
	}
	if params.Fields != "" {
		db = db.Select(PGFormatSelect(params.Fields))
	}
	if params.Order != "" {
		db = db.Order(PGFormatOrder(params.Order))
	}
	argsLen := len(args)
	if argsLen != 0 {
		if argsLen == 1 {
			db = db.Where(args[0])
		} else {
			db = db.Where(args[0], args[1:]...)
		}
	}
	return db
}

// PGCount pg count
func PGCount(model interface{}, args ...interface{}) (count int, err error) {
	db := pgClient.Model(model)
	if len(args) > 1 {
		db = db.Where(args[0], args[1:]...)
	} else if len(args) == 1 {
		db = db.Where(args[0])
	} else {
		scope := db.NewScope(model)
		db = pgClient.Table(scope.TableName())
	}
	err = db.Count(&count).Error
	return
}
