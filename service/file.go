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

package service

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v6"
	"github.com/vicanso/origin/config"
)

type (
	FileSrv      struct{}
	UploadParams struct {
		Bucket string
		Name   string
		Reader io.Reader
		Size   int64
		Opts   minio.PutObjectOptions
	}
)

var (
	minioClient *minio.Client
)

func init() {
	minioConfig := config.GetMinioConfig()
	c, err := minio.New(minioConfig.Endpoint, minioConfig.AccessKeyID, minioConfig.SecretAccessKey, minioConfig.SSL)
	if err != nil {
		panic(err)
	}
	minioClient = c
}

// Upload upload file
func (srv *FileSrv) Upload(params UploadParams) (n int64, err error) {
	return minioClient.PutObject(params.Bucket, params.Name, params.Reader, params.Size, params.Opts)
}

// Get get file
func (srv *FileSrv) Get(bucket, filename string) (*minio.Object, error) {
	return minioClient.GetObject(bucket, filename, minio.GetObjectOptions{})
}

// GetData get file data
func (srv *FileSrv) GetData(bucket, filename string) (data []byte, header http.Header, err error) {
	object, err := srv.Get(bucket, filename)
	if err != nil {
		return
	}
	statsInfo, err := object.Stat()
	if err != nil {
		return
	}
	header = make(http.Header)
	header.Set("ETag", statsInfo.ETag)
	header.Set("Content-Type", statsInfo.ContentType)

	data, err = ioutil.ReadAll(object)
	if err != nil {
		return
	}
	return
}

func (srv *FileSrv) OptimImage(reader io.Reader, imageType string, width, height int) (buffer *bytes.Buffer, err error) {
	var img image.Image
	switch imageType {
	default:
		img, _, err = image.Decode(reader)
	case "png":
		img, err = png.Decode(reader)
	case "jpg":
		fallthrough
	case "jpeg":
		img, err = jpeg.Decode(reader)
	}
	if err != nil {
		return
	}
	img = imaging.Resize(img, width, height, imaging.Lanczos)
	buffer = new(bytes.Buffer)
	switch imageType {
	default:
		err = jpeg.Encode(buffer, img, &jpeg.Options{
			Quality: 80,
		})
	case "png":
		err = png.Encode(buffer, img)
	}
	if err != nil {
		return
	}
	return
}

// func imageDecode(buf []byte, sourceType EncodeType) (img image.Image, err error) {
// 	reader := bytes.NewReader(buf)
// 	switch sourceType {
// 	default:
// 		img, _, err = image.Decode(reader)
// 	case EncodeTypeWEBP:
// 		img, err = WebpDecode(reader)
// 	case EncodeTypePNG:
// 		img, err = png.Decode(reader)
// 	case EncodeTypeJPEG:
// 		img, err = jpeg.Decode(reader)
// 	}
// 	return
// }
