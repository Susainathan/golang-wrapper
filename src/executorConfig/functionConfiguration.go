package executorConfig

import "go/xmlc-wrapper/src/serviceExecutors"

var Config = map[string]interface{}{
	"SQCExecutor": serviceExecutors.SQCExecutor,
}
