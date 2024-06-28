package filehandler

import (
	"bytes"
	"fmt"
	"go/xmlc-wrapper/src/helpers"
	"go/xmlc-wrapper/src/structs"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func GetFile(s3ClientStruct structs.S3ClientStruct, fileHandler structs.FileHandlerDetailsStruct, awsCredential structs.AwsS3AccessStruct, message map[string]string) (string, error) {
	bucket := fileHandler.FileHandlerParams.BucketName
	inputS3File := message["filePath"] + message["subFolderPath"] + message["fileName"]
	pathToLocal := fileHandler.FileHandlerParams.LocalPath + message["subFolderPath"]

	helpers.LogMessage("INFO", fmt.Sprintf("Started downloading file from the s3 path: %q", inputS3File))

	if err := os.MkdirAll(pathToLocal, os.ModePerm); err != nil {
		helpers.LogMessage("ERROR", fmt.Sprintf("Unable to create folder in local path %q, %v", pathToLocal, err))
	}

	file, err := os.Create(pathToLocal + message["fileName"])
	if err != nil {
		helpers.LogMessage("ERROR", fmt.Sprintf("Unable to create file %q, %v", pathToLocal, err))
		return "", err
	}
	defer file.Close()

	numBytes, err := s3ClientStruct.S3DownloadClient.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(inputS3File),
		})
	if err != nil {
		helpers.LogMessage("ERROR", fmt.Sprintf("Unable to download file from the s3 path: %q, %v", inputS3File, err))
		return "", err
	}
	helpers.LogMessage("INFO", fmt.Sprintf("File downloaded from the s3 path: %q, with the size of %v bytes", file.Name(), numBytes))

	return pathToLocal + message["fileName"], nil
}

func UploadFile(s3ClientStruct structs.S3ClientStruct, fileHandler structs.FileHandlerDetailsStruct, filePath string, message map[string]string) error {
	bucket := fileHandler.FileHandlerParams.BucketName
	fileName := strings.Split(filePath, "/")

	destinationPath := message["filePath"] + message["subFolderPath"] + "/" + fileName[len(fileName)-1]

	file, err := os.Open(filePath)
	if err != nil {
		helpers.LogMessage("ERROR", fmt.Sprintf("Unable to open file: %q, %v", destinationPath, err))
		return err
	}
	defer file.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		helpers.LogMessage("ERROR", fmt.Sprintf("Error reading file from the path: %q, %v", destinationPath, err))
		return err
	}

	_, err = s3ClientStruct.S3UploaderClient.Upload(&s3manager.UploadInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(destinationPath),
		Body:                 bytes.NewReader(buf.Bytes()),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		helpers.LogMessage("ERROR", fmt.Sprintf("Unable to upload file to s3 path: %q, %v", destinationPath, err))
		return err
	}

	helpers.LogMessage("INFO", fmt.Sprintf("File uploaded successfully!!! %q", destinationPath))

	return nil
}
