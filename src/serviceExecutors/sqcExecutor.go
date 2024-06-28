package serviceExecutors

import (
	"bytes"
	"errors"
	"fmt"
	"go/xmlc-wrapper/src/filehandler"
	"go/xmlc-wrapper/src/helpers"
	"os"
	"os/exec"

	"go/xmlc-wrapper/src/structs"
)

func isOutputFileExits(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func SQCExecutor(s3Client structs.S3ClientStruct, appConfig structs.SubscribedAppStruct, message map[string]string) bool {
	helpers.LogMessage("INFO", fmt.Sprintf("Started processing the request for the project token %v", message["projectToken"]))
	inputFilePath, err := filehandler.GetFile(s3Client, appConfig.FileHandlerDetails, appConfig.AwsS3Access, message)

	if err != nil {
		return false
	}

	cmd := exec.Command("docker", "run", "--rm", "-i", "-v", "/tnq/tools/mliflow:/tnq/tools/mliflow", "-v", "/tnq/data/xml-central:/tnq/data/xml-central", "-v", "/etc/ssl/certs:/etc/ssl/certs", "--name", "container_name", "mliflow_worker:1.1", "perl", "/tnq/tools/mliflow/tud2tudplus/perl/tud2tudplus.pl", fmt.Sprintf("--in_path=%s", inputFilePath), "--config_path=/tnq/tools/mliflow/tud2tudplus/conf/tud.cfg", "--process_type=tudmovement", "--server_location=alpha-vsi")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	helpers.LogMessage("INFO", fmt.Sprintf("TUD organizer core started processing the request for the project token %v", message["projectToken"]))
	err = cmd.Run()
	fmt.Println(out.String())
	if err != nil {
		fmt.Printf("%s: %s", err, stderr.String())
	}

	if !isOutputFileExits(inputFilePath) {
		helpers.LogMessage("ERROR", fmt.Sprintf("File not generated the service tud organizer for the project token %v", message["projectToken"]))
		return false
	}

	uploadeErr := filehandler.UploadFile(s3Client, appConfig.FileHandlerDetails, inputFilePath, message)

	if uploadeErr != nil {
		helpers.LogMessage("ERROR", fmt.Sprintf("Unable to complete the TUD organizer request for the project token %v", message["projectToken"]))
		return false
	}

	helpers.LogMessage("INFO", fmt.Sprintf("Completed the TUD organizer request for the project token %v", message["projectToken"]))
	return true
}
