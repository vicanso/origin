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

package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *zap.Logger
)

type (
	pgLogger struct {
	}
)

func (l *pgLogger) Printf(layout string, v ...interface{}) {
	// TODO 如果信息中带有error字段，则输出告警
	// "msg":"pg log","message":"error/Users/xieshuzhou/github/origin/service/receiver.go:33pq: column \"user_id\" contains null values"
	Default().Info("pg log",
		zap.String("message", fmt.Sprintf(layout, v...)),
	)
}

func init() {
	c := zap.NewProductionConfig()
	c.DisableCaller = true
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 只针对panic 以上的日志增加stack trace
	l, err := c.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		panic(err)
	}
	defaultLogger = l
}

// Default get default logger
func Default() *zap.Logger {
	return defaultLogger
}

// PGLogger pg logger
func PGLogger() *pgLogger {
	return new(pgLogger)
}
