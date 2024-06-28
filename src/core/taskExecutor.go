package core

import (
	"fmt"
	"go/xmlc-wrapper/src/helpers"
	"go/xmlc-wrapper/src/structs"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func ProcesseMessage(message map[string]string, sqsClient *sqs.SQS, semaphore chan struct{}, appConfig structs.SubscribedAppStruct, functionConf map[string]interface{}) {
	s3Client := getS3Client(appConfig.AwsS3Access)
	functionName := appConfig.CoreApp.CoreExecutorFunc
	// var mapMessage = make(map[string]string)

	// json.Unmarshal([]byte(message), &mapMessage)

	if fn, ok := functionConf[functionName]; ok {
		args := []reflect.Value{reflect.ValueOf(s3Client), reflect.ValueOf(appConfig), reflect.ValueOf(message)}
		result := reflect.ValueOf(fn).Call(args)
		helpers.LogMessage("INFO", fmt.Sprintf("Result: %v", result))
	} else {
		helpers.LogMessage("ERROR", fmt.Sprintf("Function not found: %v", functionName))
	}

	// _, err := sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
	// 	QueueUrl:      aws.String(*message.Attributes["QueueUrl"]),
	// 	ReceiptHandle: message.ReceiptHandle,
	// })
	// if err != nil {
	// 	fmt.Println("Error deleting message:", err)
	// }
	<-semaphore
}

func getS3Client(s3Config structs.AwsS3AccessStruct) structs.S3ClientStruct {

	var s3Struct structs.S3ClientStruct

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s3Config.Region),
		Credentials: credentials.NewStaticCredentials(s3Config.AccessKey, s3Config.SecretKey, ""),
	})

	if err != nil {
		fmt.Println("Error creating session:", err)
		s3Struct.Error = true
		return s3Struct
	}

	s3Struct.S3DownloadClient = s3manager.NewDownloader(sess)
	s3Struct.S3UploaderClient = s3manager.NewUploader(sess)
	return s3Struct
}
