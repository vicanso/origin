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
	qrcode "github.com/skip2/go-qrcode"
)

type (
	QRCodeInfo struct {
		Data []byte `json:"data,omitempty"`
		Size int    `json:"size,omitempty"`
	}
)

// GetQRCode get a qrcode image
func GetQRCode(value string, size int) (info *QRCodeInfo, err error) {
	buf, err := qrcode.Encode(value, qrcode.High, size)
	if err != nil {
		return
	}
	info = &QRCodeInfo{
		Data: buf,
		Size: size,
	}
	return
}
