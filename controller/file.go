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
	"bytes"
	"net/http"
	"strings"

	"github.com/vicanso/origin/cs"

	"github.com/minio/minio-go/v6"
	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	"github.com/vicanso/origin/router"
	"github.com/vicanso/origin/service"
	"github.com/vicanso/origin/util"
	"github.com/vicanso/origin/validate"
)

type (
	fileCtrl struct{}
)

type (
	fileUploadParams struct {
		Bucket string `json:"bucket,omitempty" validate:"xFileBucket"`
	}
)

const (
	fileCategory    = "file"
	filePreviwRoute = "/v1/preview/{bucket}/{filename}"
)

var validContentTypes = []string{
	"image/jpeg",
	"image/png",
}

var invalidContentType = &hes.Error{
	StatusCode: http.StatusBadRequest,
	Message:    "不支持该文件类型",
	Category:   fileCategory,
}

func init() {
	ctrl := fileCtrl{}
	g := router.NewGroup("/files")
	g.POST(
		"/v1",
		loadUserSession,
		newCheckGroupsMiddleware([]string{
			cs.UserGroupMarketing,
		}),
		newTracker(cs.ActionFileUpload),
		ctrl.upload,
	)

	g.GET(filePreviwRoute, ctrl.preview)
}

// upload file upload
func (ctrl fileCtrl) upload(c *elton.Context) (err error) {

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return
	}
	contentType := header.Header.Get("Content-Type")
	if !util.ContainsString(validContentTypes, contentType) {
		err = invalidContentType
		return
	}
	params := fileUploadParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	us := getUserSession(c)
	filename := util.GenUlid() + "." + strings.Split(contentType, "/")[1]
	_, err = fileSrv.Upload(service.UploadParams{
		Bucket: params.Bucket,
		Name:   filename,
		Reader: file,
		Size:   header.Size,
		Opts: minio.PutObjectOptions{
			ContentType: contentType,
			UserTags: map[string]string{
				"account": us.GetAccount(),
			},
		},
	})
	if err != nil {
		return
	}
	url := strings.Replace("/files"+filePreviwRoute, "{bucket}", params.Bucket, 1)
	url = strings.Replace(url, "{filename}", filename, 1)
	c.Body = map[string]string{
		"url": url,
	}
	return
}

// preview file preview
func (ctrl fileCtrl) preview(c *elton.Context) (err error) {
	bucket := c.Param("bucket")
	filename := c.Param("filename")
	data, header, err := fileSrv.GetData(bucket, filename)
	if err != nil {
		return
	}
	// 客户端缓存一周，缓存服务器缓存10分钟
	c.CacheMaxAge("168h", "10m")
	for k, values := range header {
		for _, v := range values {
			c.AddHeader(k, v)
		}
	}
	c.BodyBuffer = bytes.NewBuffer(data)
	return
}
