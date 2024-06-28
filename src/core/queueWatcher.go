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
		message := fmt.Sprintf("Error creating session: %v", err)
		helpers.LogMessage("ERROR", message)
		return
	}

	serviceClient := sqs.New(sess)

	message := fmt.Sprintf("Initiating SQS to watch message for the app %v...", appConfig.App)
	helpers.LogMessage("INFO", message)

	for {
		select {
		case <-initChan:
			message := fmt.Sprintf("Stopping goroutine for the app %v", appConfig.App)
			helpers.LogMessage("INFO", message)
			return
		case semaphore <- struct{}{}:
			result, err := serviceClient.ReceiveMessage(&sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(appConfig.QueueDetails.QueueParams.Url),
				MaxNumberOfMessages: aws.Int64(1),
				VisibilityTimeout:   aws.Int64(int64(appConfig.CoreApp.ProcessingTimeout) + 60),
				WaitTimeSeconds:     aws.Int64(1),
			})
			if err != nil {
				message := fmt.Sprintf("Error receiving message: %v", err)
				helpers.LogMessage("ERROR", message)
				<-semaphore
				continue
			}

			message := map[string]string{
				"projectToken":  "1c4b56ba-19cd-44c0-a7b2-43441314ef1c",
				"filePath":      "/xml-central/sqc/1c4b56ba-19cd-44c0-a7b2-43441314ef1c",
				"subFolderPath": "/f830a848-7d48-4e24-9d50-a37d199fa86c",
				"fileName":      "/_tud.xml",
			}

			if len(result.Messages) > 0 {
				<-semaphore
				continue
			} else {
				fucntionConf := executorConfig.Config
				helpers.LogMessage("INFO", fmt.Sprintf("Received message for the app %v...", appConfig.App))

				// go ProcesseMessage(*result.Messages[0], serviceClient, semaphore, appConfig, fucntionConf)
				go ProcesseMessage(message, serviceClient, semaphore, appConfig, fucntionConf)

			}
		}
	}
}
