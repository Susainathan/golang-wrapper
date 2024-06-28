package helpers

import (
	"go/xmlc-wrapper/src/structs"
)

func isLogAvailable(value string, element [6]string) bool {
	for _, i := range element {
		if i == value {
			return true
		}
	}
	return false
}

func subscribedAppValidation(WorkerConfig []structs.SubscribedAppStruct, accessConfig map[string]map[string]string) []structs.SubscribedAppStruct {
	apps := []structs.SubscribedAppStruct{}

	for _, app := range WorkerConfig {
		if app.Enabled {
			var queueConfigStruct structs.AwsSqsAccessStruct
			var apiConfigStruct structs.ApiAccessStruct
			var s3ConfigStruct structs.AwsS3AccessStruct

			queueConfigData := accessConfig[app.QueueDetails.Access]
			apiConfigData := accessConfig[app.Webhook.Access]
			s3ConfigData := accessConfig[app.FileHandlerDetails.Access]

			if queueConfigData != nil {
				queueConfigStruct.AccessKey = queueConfigData["accessKey"]
				queueConfigStruct.SecretKey = queueConfigData["secretKey"]
				queueConfigStruct.Region = queueConfigData["region"]
			}
			if apiConfigData != nil {
				apiConfigStruct.AccessKey = apiConfigData["accessName"]
				apiConfigStruct.AccessName = apiConfigData["accessKey"]
			}
			if s3ConfigData != nil {
				s3ConfigStruct.AccessKey = s3ConfigData["accessKey"]
				s3ConfigStruct.SecretKey = s3ConfigData["secretKey"]
				s3ConfigStruct.Region = s3ConfigData["region"]
			}

			app.AwsSqsAccess = queueConfigStruct
			app.AwsS3Access = s3ConfigStruct
			app.ApiAccess = apiConfigStruct

			apps = append(apps, app)
		}
	}
	return apps
}

func ValidateConfig(WorkerConfig structs.WorkerConfigStruct) []structs.SubscribedAppStruct {

	config := WorkerConfig

	logLevels := [6]string{"DEBUG", "INFO", "WARNING", "WARN", "ERROR", "CRITICAL"}

	if !config.Enabled {
		panic("Worker not enabled..")
	}
	if !isLogAvailable(config.LogLevel, logLevels) {
		panic("Invalid worker config.")
	}
	if len(WorkerConfig.Accesses) == 0 {
		panic("Invalid access config.")
	}
	subscribedApps := subscribedAppValidation(config.Subscribed, WorkerConfig.Accesses)

	if len(subscribedApps) == 0 {
		panic("No subscribed app found.")
	}

	return subscribedApps
}
