package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"net/http"
	"regexp"
	"strings"
)

type anthropicService struct {
	cfg *config.Config
}

type AnthropicService interface {
	FetchSummary(emailHtml string) string
	FetchActionItems(emailHtml string) []string
}

func (s *anthropicService) FetchSummary(emailHtml string) string {
	reqBody := map[string]interface{}{
		"prompt": s.cfg.Anthropic.SummaryPrompt + emailHtml,
		"model":  "claude-2",
	}

	jsonBody, _ := json.Marshal(reqBody)
	reqReader := bytes.NewReader(jsonBody)

	req, _ := http.NewRequest("POST", s.cfg.Anthropic.ApiPath+"/ask", reqReader)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Openline-API-KEY", s.cfg.Anthropic.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	// Print summarized email
	var data map[string]string
	json.NewDecoder(resp.Body).Decode(&data)
	return extractSummaryText(data["completion"])
}

func (s *anthropicService) FetchActionItems(emailHtml string) []string {
	reqBody := map[string]interface{}{
		"prompt": s.cfg.Anthropic.ActionItemsPromp + emailHtml,
		"model":  "claude-2",
	}

	jsonBody, _ := json.Marshal(reqBody)
	reqReader := bytes.NewReader(jsonBody)

	req, _ := http.NewRequest("POST", s.cfg.Anthropic.ApiPath, reqReader)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Openline-API-KEY", s.cfg.Anthropic.ApiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	// Print summarized email
	var data map[string]string
	json.NewDecoder(resp.Body).Decode(&data)

	return extractActionPoints(data["completion"])
}

func extractSummaryText(input string) string {
	colonIndex := strings.Index(input, ":")
	if colonIndex == -1 {
		return "" // Return an empty string if no colon is found
	}

	summary := strings.TrimSpace(input[colonIndex+1:])
	return summary
}

func extractActionPoints(text string) []string {
	// Regular expression pattern to match bullet points
	pattern := `\s*-\s*(.*)`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Split the text into lines
	lines := strings.Split(text, "\n")

	// Store the action points
	actionPoints := []string{}

	// Iterate over each line and extract the action points
	for _, line := range lines {
		// Check if the line matches the bullet point pattern
		matches := re.FindStringSubmatch(line)
		if len(matches) == 2 {
			actionPoints = append(actionPoints, matches[1])
		}
	}

	return actionPoints
}

func NewAnthropicService(cfg *config.Config) AnthropicService {
	return &anthropicService{
		cfg: cfg,
	}
}
