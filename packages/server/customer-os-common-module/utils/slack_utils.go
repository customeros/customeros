package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"io"
	"net/http"
)

func SendSlackMessage(c context.Context, slackWehbookUrl, text string) error {
	span, _ := opentracing.StartSpanFromContext(c, "SlackUtils.SendSlackMessage")
	defer span.Finish()

	// Create a struct to hold the JSON data
	type SlackMessage struct {
		Text string `json:"text"`
	}
	message := SlackMessage{Text: text}

	// Convert struct to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Send POST request
	resp, err := http.Post(slackWehbookUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(log.String("response.body", string(responseBody)))

	return nil
}
