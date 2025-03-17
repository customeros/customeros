package data_fields

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type IndustryCode struct {
	Code       string  `json:"code"`
	Confidence float64 `json:"confidence"`
}

type IndustryCodeResponse struct {
	Industry IndustryCode `json:"industry"`
}

type IndustryCodeResponseValidator struct {
	ValidateFormat bool     // Whether to validate NAICS code format
	DisallowEmpty  bool     // Whether to disallow empty code
	BlockedCodes   []string // NAICS codes that shouldn't appear
	MinConfidence  float64  // Minimum confidence score allowed
	MaxConfidence  float64  // Maximum confidence score allowed
}

func NewIndustryCodeResponseValidator() *IndustryCodeResponseValidator {
	return &IndustryCodeResponseValidator{
		ValidateFormat: true,
		DisallowEmpty:  true,
		BlockedCodes:   []string{},
		MinConfidence:  0.0,
		MaxConfidence:  1.0,
	}
}

func (v *IndustryCodeResponseValidator) IsValidJSON(jsonStr string) bool {
	var response IndustryCodeResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	return err == nil
}

func (v *IndustryCodeResponseValidator) ValidateResponse(responseStr string) (*IndustryCode, error) {
	// Parse the JSON
	var response IndustryCodeResponse
	err := json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	// Clean the NAICS code
	cleaned := strings.TrimSpace(response.Industry.Code)

	// Check if empty
	if v.DisallowEmpty && cleaned == "" {
		return nil, fmt.Errorf("empty NAICS code")
	}

	// Validate NAICS code format (6-digit number)
	if v.ValidateFormat {
		match, _ := regexp.MatchString(`^\d{6}$`, cleaned)
		if !match {
			return nil, fmt.Errorf("invalid NAICS code format: must be a 6-digit number")
		}
	}

	// Check for blocked codes
	for _, code := range v.BlockedCodes {
		if cleaned == code {
			return nil, fmt.Errorf("blocked NAICS code: %s", code)
		}
	}

	// Validate confidence score
	if response.Industry.Confidence < v.MinConfidence {
		return nil, fmt.Errorf("confidence score %.2f is below minimum %.2f",
			response.Industry.Confidence, v.MinConfidence)
	}

	if response.Industry.Confidence > v.MaxConfidence {
		return nil, fmt.Errorf("confidence score %.2f is above maximum %.2f",
			response.Industry.Confidence, v.MaxConfidence)
	}

	// Return validated result
	return &IndustryCode{
		Code:       cleaned,
		Confidence: response.Industry.Confidence,
	}, nil
}

func (v *IndustryCodeResponseValidator) GetExpectedSchema() string {
	return `{
  "industry": {
    "code": "541330",
    "confidence": 0.85
  }
}`
}
