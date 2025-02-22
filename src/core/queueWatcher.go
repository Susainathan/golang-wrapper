package core

import (
	"fmt"
	"go/xmlc-wrapper/src/executorConfig"
	"go/xmlc-wrapper/src/helpers"
	"go/xmlc-wrapper/src/structs"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func QueueWatcher(appConfig structs.SubscribedAppStruct, wg *sync.WaitGroup, initChan chan struct{}, semaphore chan struct{}) {
	defer wg.Done()

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(appConfig.AwsSqsAccess.Region),
		Credentials: credentials.NewStaticCredentials(appConfig.AwsSqsAccess.AccessKey, appConfig.AwsSqsAccess.SecretKey, ""),
	})

	if err != nil {
		helpers.LogMessage("ERROR", fmt.Sprintf("Error creating session: %v", err))
		return
	}

	serviceClient := sqs.New(sess)
	helpers.LogMessage("INFO", fmt.Sprintf("Initiating SQS to watch messages for the app %v...", appConfig.App))

	for {
		select {
		case <-initChan:
			helpers.LogMessage("INFO", fmt.Sprintf("Stopping goroutine for the app %v", appConfig.App))
			return
		case semaphore <- struct{}{}:
			// Receive messages from SQS
			result, err := serviceClient.ReceiveMessage(&sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(appConfig.QueueDetails.QueueParams.Url),
				MaxNumberOfMessages: aws.Int64(1), // Fetch 1 message at a time
				VisibilityTimeout:   aws.Int64(int64(appConfig.CoreApp.ProcessingTimeout) + 60),
				WaitTimeSeconds:     aws.Int64(1),
			})
			if err != nil {
				helpers.LogMessage("ERROR", fmt.Sprintf("Error receiving message: %v", err))
				<-semaphore // Release semaphore slot on error
				continue
			}

			if len(result.Messages) == 0 {
				<-semaphore
				continue
			}

			go func(msg *sqs.Message) {
				defer func() {
					<-semaphore
				}()

				fucntionConf := executorConfig.Config
				helpers.LogMessage("INFO", fmt.Sprintf("Processing message for app %v...", appConfig.App))

				ProcesseMessage(msg, serviceClient, appConfig, fucntionConf)

				_, err := serviceClient.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      aws.String(appConfig.QueueDetails.QueueParams.Url),
					ReceiptHandle: msg.ReceiptHandle,
				})
				if err != nil {
					helpers.LogMessage("ERROR", fmt.Sprintf("Error deleting message: %v", err))
				} else {
					helpers.LogMessage("INFO", fmt.Sprintf("Message processed and deleted for app %v", appConfig.App))
				}
			}(result.Messages[0])
		}
	}
}
