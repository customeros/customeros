package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"io"
	"net/http"
)

func SendSlackMessage(c context.Context, slackWehbookUrl, text string) error {
	spans, _ := telemetry.StartServiceSpan(c, "SlackUtils.SendSlackMessage")
	defer spans.Finish()

	// Create a struct to hold the JSON data
	type SlackMessage struct {
		Text string `json:"text"`
	}
	message := SlackMessage{Text: text}

	// Convert struct to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	// Send POST request
	resp, err := http.Post(slackWehbookUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		spans.TraceError(err)
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	spans.LogKV("response.body", string(responseBody))

	return nil
}
