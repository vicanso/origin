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

package config

import (
	"bytes"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
	"github.com/vicanso/origin/validate"
)

var (
	box       = packr.New("config", "../configs")
	env       = os.Getenv("GO_ENV")
	appName   string
	version   string
	buildedAt string
)

const (
	ApplicationStatusStopped int32 = iota
	ApplicationStatusRunning
	ApplicationStatusStopping
)

type (
	// RedisOptions redis options
	RedisOptions struct {
		Addr     string `validate:"min=5,max=30"`
		Password string
		DB       int
		// 慢请求时长
		Slow time.Duration
		// 最大的正在处理请求量
		MaxProcessing uint32
	}
	// SessionConfig session's config
	SessionConfig struct {
		TTL        time.Duration
		Key        string `validate:"ascii"`
		CookiePath string `validate:"ascii"`
	}
	// MailConfig mail's config
	MailConfig struct {
		Host     string `validate:"hostname"`
		Port     int    `validate:"number"`
		User     string `validate:"email"`
		Password string `validate:"min=1,max=100"`
	}

	// Influxdb config
	InfluxdbConfig struct {
		Bucket        string `validate:"min=1,max=50"`
		Org           string `validate:"min=1,max=100"`
		URI           string `validate:"url"`
		Token         string `validate:"ascii"`
		BatchSize     uint   `validate:"min=1,max=5000"`
		FlushInterval time.Duration
		Disabled      bool
	}

	// PostgresConfig postgres config
	PostgresConfig struct {
		Slow                time.Duration
		MaxQueryProcessing  uint32
		MaxUpdateProcessing uint32
		DisableAutoMigrate  bool
	}

	// MinioConfig minio config
	MinioConfig struct {
		Endpoint        string `validate:"hostname_port"`
		AccessKeyID     string `validate:"min=3"`
		SecretAccessKey string `validate:"min=6"`
		SSL             bool
	}

	// TinyConfig tiny config
	TinyConfig struct {
		Host string `validate:"ip"`
		Port int    `validate:"number"`
	}
)

const (
	// Dev development env
	Dev = "dev"
	// Test test env
	Test = "test"
	// Production production env
	Production = "production"

	defaultListen   = ":7001"
	defaultTrackKey = "jt"
)

var (
	defaultViper = viper.New()

	applicationStatus = ApplicationStatusStopped
)

func init() {
	configType := "yml"
	v := viper.New()
	defaultViper.SetConfigType(configType)
	v.SetConfigType(configType)

	configExt := "." + configType
	data, err := box.Find("default" + configExt)
	if err != nil {
		panic(err)
	}
	// 读取默认配置中的所有配置
	err = v.ReadConfig(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	configs := v.AllSettings()
	// 将default中的配置全部以默认配置写入
	for k, v := range configs {
		defaultViper.SetDefault(k, v)
	}

	// 根据当前运行环境配置读取
	envConfigFile := GetENV() + configExt
	data, err = box.Find(envConfigFile)
	if err != nil {
		panic(err)
	}
	// 读取当前运行环境对应的配置
	err = defaultViper.ReadConfig(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	appName = GetString("app")
}

func validatePanic(v interface{}) {
	err := validate.Do(v, nil)
	if err != nil {
		panic(err)
	}
}

// SetApplicationStatus set application status
func SetApplicationStatus(status int32) {
	atomic.StoreInt32(&applicationStatus, status)
}

// ApplicationIsRunning application is running
func ApplicationIsRunning() bool {
	return atomic.LoadInt32(&applicationStatus) == ApplicationStatusRunning
}

func GetAppName() string {
	return appName
}

// GetVersion get application's version
func GetVersion() string {
	return version
}

// SetVersion set application's version
func SetVersion(v string) {
	version = v
}

// GetBuildedAt get the builded time of application
func GetBuildedAt() string {
	return buildedAt
}

// SetBuildedAt set the builded time of application
func SetBuildedAt(v string) {
	buildedAt = v
}

// GetENV get go env
func GetENV() string {
	if env == "" {
		return Dev
	}
	return env
}

// GetRequestLimit get request limit
func GetRequestLimit() uint32 {
	return GetUint32("requestLimit")
}

// GetBool viper get bool
func GetBool(key string) bool {
	return defaultViper.GetBool(key)
}

// GetInt viper get int
func GetInt(key string) int {
	return defaultViper.GetInt(key)
}

// GetUint viper get uint
func GetUint(key string) uint {
	return defaultViper.GetUint(key)
}

// GetUint32 viper get uint32
func GetUint32(key string) uint32 {
	return defaultViper.GetUint32(key)
}

// GetIntDefault get int with default value
func GetIntDefault(key string, defaultValue int) int {
	v := GetInt(key)
	if v != 0 {
		return v
	}
	return defaultValue
}

// GetUint32Default get uint32 with default value
func GetUint32Default(key string, defaultValue uint32) uint32 {
	v := GetUint32(key)
	if v != 0 {
		return v
	}
	return defaultValue
}

// GetString viper get string
func GetString(key string) string {
	return defaultViper.GetString(key)
}

// GetStringFromENV get string from env, if not exists, it will return the value of config
func GetStringFromENV(key string) string {
	value := GetString(key)
	v := os.Getenv(value)
	if v != "" {
		return v
	}
	return value
}

// GetStringDefault get string with default value
func GetStringDefault(key, defaultValue string) string {
	v := GetString(key)
	if v != "" {
		return v
	}
	return defaultValue
}

// GetDuration viper get duration
func GetDuration(key string) time.Duration {
	return defaultViper.GetDuration(key)
}

// GetDurationDefault get duration with default value
func GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	v := GetDuration(key)
	if v != 0 {
		return v
	}
	return defaultValue
}

// GetStringSlice viper get string slice
func GetStringSlice(key string) []string {
	return defaultViper.GetStringSlice(key)
}

// GetStringMap get string map
func GetStringMap(key string) map[string]interface{} {
	return defaultViper.GetStringMap(key)
}

// GetListen get server listen address
func GetListen() string {
	return GetStringDefault("listen", defaultListen)
}

// GetTrackKey get the track cookie key
func GetTrackKey() string {
	return GetStringDefault("track", defaultTrackKey)
}

// GetRedisConfig get redis config
func GetRedisConfig() (options RedisOptions, err error) {
	prefix := "redis."
	uri := GetStringFromENV(prefix + "uri")
	defaultSlow := 300 * time.Millisecond
	defaultMaxProcessing := 1000
	if len(uri) != 0 {
		urlInfo, e := url.Parse(uri)
		if e != nil {
			err = e
			return
		}
		pass, _ := urlInfo.User.Password()
		query := urlInfo.Query()
		db, _ := strconv.Atoi(query.Get("db"))
		slow, _ := time.ParseDuration(query.Get("slow"))
		if slow == 0 {
			slow = defaultSlow
		}
		maxProcessing, _ := strconv.Atoi(query.Get("maxProcessing"))
		if maxProcessing == 0 {
			maxProcessing = defaultMaxProcessing
		}

		options = RedisOptions{
			Addr:          urlInfo.Host,
			Password:      pass,
			DB:            db,
			Slow:          slow,
			MaxProcessing: uint32(maxProcessing),
		}
	} else {
		options = RedisOptions{
			Addr:          GetStringFromENV(prefix + "addr"),
			Password:      GetStringFromENV(prefix + "password"),
			DB:            GetInt(prefix + "db"),
			Slow:          GetDurationDefault(prefix+"slow", defaultSlow),
			MaxProcessing: GetUint32Default(prefix+"maxProcessing", uint32(defaultMaxProcessing)),
		}
	}
	validatePanic(&options)
	return
}

// GetPostgresConnectString get postgres connect string
func GetPostgresConnectString() string {
	keys := []string{
		"host",
		"port",
		"user",
		"dbname",
		"password",
		"sslmode",
	}
	arr := []string{}
	prefix := "postgres."
	uri := GetStringFromENV(prefix + "uri")
	if uri != "" {
		return uri
	}
	for _, key := range keys {
		k := prefix + key
		value := GetString(k)
		// 密码与用户名支持env中获取
		if key == "password" || key == "user" {
			value = GetStringFromENV(k)
		}
		if value != "" {
			arr = append(arr, key+"="+value)
		}
	}
	return strings.Join(arr, " ")
}

// GetPostgresConfig get postgres config
func GetPostgresConfig() PostgresConfig {
	prefix := "postgres."
	slow := GetDuration(prefix + "slow")
	if slow == 0 {
		slow = time.Second
	}
	maxQueryProcessing := GetUint32Default(prefix+"maxQueryProcessing", 1000)
	maxUpdateProcessing := GetUint32Default(prefix+"maxUpdateProcessing", 500)

	return PostgresConfig{
		Slow:                slow,
		MaxQueryProcessing:  maxQueryProcessing,
		MaxUpdateProcessing: maxUpdateProcessing,
		DisableAutoMigrate:  GetBool(prefix + "disableAutoMigrate"),
	}
}

// GetSessionConfig get sesion config
func GetSessionConfig() SessionConfig {
	prefix := "session."
	sessConfig := SessionConfig{
		TTL:        GetDuration(prefix + "ttl"),
		Key:        GetString(prefix + "key"),
		CookiePath: GetString(prefix + "path"),
	}
	// 如果session设置过短，则使用默认为24小时
	if sessConfig.TTL < time.Second {
		sessConfig.TTL = 24 * time.Hour
	}
	validatePanic(&sessConfig)
	return sessConfig
}

// GetSignedKeys get signed keys
func GetSignedKeys() []string {
	return GetStringSlice("keys")
}

// GetMailConfig get mail config
func GetMailConfig() MailConfig {
	prefix := "mail."
	mailConfig := MailConfig{
		Host:     GetString(prefix + "host"),
		Port:     GetInt(prefix + "port"),
		User:     GetStringFromENV(prefix + "user"),
		Password: GetStringFromENV(prefix + "password"),
	}
	validatePanic(&mailConfig)
	return mailConfig
}

// GetInfluxdbConfig get influxdb config
func GetInfluxdbConfig() InfluxdbConfig {
	prefix := "influxdb."
	influxdbConfig := InfluxdbConfig{
		URI:           GetStringFromENV(prefix + "uri"),
		Bucket:        GetString(prefix + "bucket"),
		Org:           GetString(prefix + "org"),
		Token:         GetStringFromENV(prefix + "token"),
		BatchSize:     GetUint(prefix + "batchSize"),
		FlushInterval: GetDuration(prefix + "flushInterval"),
		Disabled:      GetBool(prefix + "disabled"),
	}
	validatePanic(&influxdbConfig)
	return influxdbConfig
}

// GetMinioConfig get minio config
func GetMinioConfig() MinioConfig {
	prefix := "minio."
	minioConfig := MinioConfig{
		Endpoint:        GetString(prefix + "endpoint"),
		AccessKeyID:     GetStringFromENV(prefix + "accessKeyID"),
		SecretAccessKey: GetStringFromENV(prefix + "secretAccessKey"),
		SSL:             GetBool(prefix + "ssl"),
	}
	validatePanic(&minioConfig)
	return minioConfig
}

// GetTinyConfig get tiny config
func GetTinyConfig() TinyConfig {
	prefix := "tiny."
	tinyConfig := TinyConfig{
		Host: GetString(prefix + "host"),
		Port: GetInt(prefix + "port"),
	}
	validatePanic(&tinyConfig)
	return tinyConfig
}
