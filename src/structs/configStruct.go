package structs

type WorkerConfigStruct struct {
	Enabled    bool                  `json:"enabled"`
	LogLevel   string                `json:"logLevel"`
	Subscribed []SubscribedAppStruct `json:"subscribed"`
	Accesses   map[string]map[string]string
}

type SubscribedAppStruct struct {
	App                string                   `json:"app"`
	Threads            int64                    `json:"threads"`
	Enabled            bool                     `json:"enabled"`
	CoreApp            CoreAppConfigStruct      `json:"coreApp"`
	QueueDetails       QueueConfigStruct        `json:"queueDetails"`
	FileHandlerDetails FileHandlerDetailsStruct `json:"fileHandlerDetails"`
	Webhook            WebhookStruct            `json:"webhook"`
	ApiAccess          ApiAccessStruct
	AwsSqsAccess       AwsSqsAccessStruct
	AwsS3Access        AwsS3AccessStruct
}

type CoreAppConfigStruct struct {
	CoreExecutorFunc  string            `json:"coreExecutorFunc"`
	ProcessingTimeout int               `json:"processingTimeout"`
	ServiceParams     map[string]string `json:"serviceParams"`
}

type QueueConfigStruct struct {
	Service     string                  `json:"service"`
	Access      string                  `json:"access"`
	QueueParams QueueParamsConfigStruct `json:"queueParams"`
}

type QueueParamsConfigStruct struct {
	QueueType string `json:"queueType"`
	Url       string `json:"url"`
}

type FileHandlerDetailsStruct struct {
	Service           string                  `json:"service"`
	Access            string                  `json:"access"`
	FileHandlerParams FileHandlerParamsStruct `json:"params"`
}

type FileHandlerParamsStruct struct {
	BucketName     string `json:"bucketName"`
	LocalPath      string `json:"localPath"`
	OutputFileName string `json:"outputFileName"`
}

type WebhookStruct struct {
	Enabled string `json:"enabled"`
	Access  string `json:"access"`
	Api     string `json:"api"`
}

type ApiAccessStruct struct {
	AccessName string
	AccessKey  string
}

type AwsSqsAccessStruct struct {
	Region    string
	AccessKey string
	SecretKey string
}

type AwsS3AccessStruct struct {
	Region    string
	AccessKey string
	SecretKey string
}
