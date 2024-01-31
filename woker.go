package main

import (
	"bytes"
	"cl-inter/helpers"
	"encoding/json"
	"log/slog"
	"net/http"
)

func worker(requestCh <-chan map[string]interface{}) {
	for data := range requestCh {
		outputMap := helpers.GenerateOutputData(data)

		outputJSON, err := json.Marshal(outputMap)
		if err != nil {
			slog.Error("Error:", err)
			return
		}

		url := "https://webhook.site/c30a1e1f-c589-4894-bb27-8780e0f87c68"

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(outputJSON))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			slog.Error("Error:", err)
			return
		}
		defer resp.Body.Close()
	}
}
