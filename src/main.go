package main

import (
	"encoding/json"
	"go/xmlc-wrapper/src/core"
	"go/xmlc-wrapper/src/helpers"
	"go/xmlc-wrapper/src/structs"
	"os"
	"path/filepath"
)

func readJson(file string, readFor string) (structs.WorkerConfigStruct, map[string]map[string]string) {
	byteValue, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	var result map[string]map[string]string
	var WorkerConfig structs.WorkerConfigStruct

	if readFor == "worker" {
		json.Unmarshal(byteValue, &WorkerConfig)
		return WorkerConfig, result
	} else {
		json.Unmarshal([]byte(byteValue), &result)
		return WorkerConfig, result
	}
}

func main() {
	helpers.InitLogger()

	WorkerConfig, _ := readJson(filepath.Base(".")+"/config/worker.json", "worker")
	_, AccessConfig := readJson(filepath.Base(".")+"/config/access.json", "access")
	// WorkerConfig, _ := readJson("/home/susainathan/Desktop/learning/go/golang-wrapper/config/worker.json", "worker")
	// _, AccessConfig := readJson("/home/susainathan/Desktop/learning/go/golang-wrapper/config/access.json", "access")
	WorkerConfig.Accesses = AccessConfig

	configData := helpers.ValidateConfig(WorkerConfig)

	core.CreateBaseChannels(configData)
}
