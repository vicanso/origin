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

package util

import (
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/jinzhu/now"
)

var (
	mockTime int64
)

// SetMockTime set mock time
func SetMockTime(v string) {
	if v == "" {
		atomic.StoreInt64(&mockTime, 0)
		return
	}
	t, _ := time.Parse(time.RFC3339, v)
	seconds := t.Unix()
	// 设置的时间有误，不调整
	if seconds < 0 {
		return
	}
	atomic.StoreInt64(&mockTime, seconds)
}

// Now get the now time
func Now() time.Time {
	// 正式环境不提供mock time
	if IsProduction() {
		return time.Now()
	}
	v := atomic.LoadInt64(&mockTime)
	if v == 0 {
		return time.Now()
	}
	return time.Unix(v, 0)
}

// NowString get the now time string of time RFC3339
func NowString() string {
	return Now().Format(time.RFC3339)
}

// UTCNow get the utc time
func UTCNow() time.Time {
	return Now().UTC()
}

// UTCNowString get the utc time string of time RFC3339
func UTCNowString() string {
	return UTCNow().Format(time.RFC3339)
}

// ParseTime parse time
func ParseTime(str string) (time.Time, error) {
	return time.Parse(time.RFC3339, str)
}

// FormatTime format time
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ChinaNow get the now time of china
func ChinaNow() (time.Time, error) {
	t := Now()
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return t, err
	}
	return t.In(loc), nil
}

// ChinaToday china today
func ChinaToday() (time.Time, error) {
	t, err := ChinaNow()
	return now.With(t).BeginningOfDay(), err
}

// ChinaYesterday china yesterday
func ChinaYesterday() (time.Time, error) {
	t, err := ChinaNow()
	if err != nil {
		return t, err
	}
	t = t.AddDate(0, 0, -1)
	return now.With(t).BeginningOfDay(), err
}

// IsBetween now is between begin and end
func IsBetween(begin *time.Time, end *time.Time) bool {
	now := Now().Unix()
	// 如果开始时间大于当前时间，未开始
	if begin != nil && begin.Unix() > now {
		return false
	}
	// 如果结束时间少于当前时间，已结束
	if end != nil && end.Unix() < now {
		return false
	}
	return true
}

// NewTimeWithRandomNS new a time with random nano seconds
func NewTimeWithRandomNS(timestamp int64) time.Time {
	rand.Seed(time.Now().UnixNano())
	sec := timestamp / 1000
	ms := timestamp % 1000
	ns := ms*10e6 + time.Now().UnixNano()%10e6
	return time.Unix(sec, ns)
}
