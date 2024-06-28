package core

import (
	"go/xmlc-wrapper/src/helpers"
	"go/xmlc-wrapper/src/structs"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func CreateBaseChannels(configData []structs.SubscribedAppStruct) {
	var wg sync.WaitGroup
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	totalConfigData := len(configData)

	stopInitChan := make(chan struct{})
	stopChannels := make([]chan struct{}, totalConfigData)

	channelCont := 0

	for i, app := range configData {
		stopChannels[i] = make(chan struct{})
		semaphore := make(chan struct{}, app.Threads)
		go QueueWatcher(app, &wg, stopInitChan, semaphore)
		channelCont++
	}

	wg.Add(totalConfigData)
	<-interrupt

	close(stopInitChan)

	for _, stopChan := range stopChannels {
		close(stopChan)
	}
	wg.Wait()

	helpers.LogMessage("INFO", "All channels are stoped.")
}
