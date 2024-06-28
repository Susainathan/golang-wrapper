package structs

import "github.com/aws/aws-sdk-go/service/s3/s3manager"

type S3ClientStruct struct {
	S3DownloadClient *s3manager.Downloader
	Error            bool
	S3UploaderClient *s3manager.Uploader
}

type BaseResponse struct {
	Message string `json:"Message"`
}

type DashboardStruct struct {
	Response string `json:"response"`
}

type LogLevel struct {
	Debug   string
	Error   string
	Info    string
	Warning string
}
