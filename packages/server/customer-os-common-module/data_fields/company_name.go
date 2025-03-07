package data_fields

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CompanyIdentification struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

type CompanyIdentificationResponse struct {
	Company CompanyIdentification `json:"company"`
}

type CompanyResponseValidator struct {
	MaxNameLen    int      // Maximum length of company name
	MinNameLen    int      // Minimum length of company name
	DisallowEmpty bool     // Whether to disallow empty company name
	BlockedWords  []string // Words that shouldn't appear in company name
	MinConfidence float64  // Minimum confidence score allowed
	MaxConfidence float64  // Maximum confidence score allowed
}

func NewCompanyResponseValidator() *CompanyResponseValidator {
	return &CompanyResponseValidator{
		MaxNameLen:    100,
		MinNameLen:    1,
		DisallowEmpty: true,
		BlockedWords:  []string{},
		MinConfidence: 0.0,
		MaxConfidence: 1.0,
	}
}

func (v *CompanyResponseValidator) IsValidJSON(jsonStr string) bool {
	var response CompanyIdentificationResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	return err == nil
}

func (v *CompanyResponseValidator) ValidateResponse(jsonStr string) (*CompanyIdentification, error) {
	// Parse the JSON
	var response CompanyIdentificationResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	// Clean the company name
	cleaned := strings.TrimSpace(response.Company.Name)

	// Check if empty
	if v.DisallowEmpty && cleaned == "" {
		return nil, fmt.Errorf("empty company name")
	}

	// Check name length
	if len(cleaned) > v.MaxNameLen {
		return nil, fmt.Errorf("company name exceeds maximum length: got %d, max is %d",
			len(cleaned), v.MaxNameLen)
	}

	if len(cleaned) < v.MinNameLen {
		return nil, fmt.Errorf("company name is too short: got %d, min is %d",
			len(cleaned), v.MinNameLen)
	}

	// Check for blocked words
	for _, word := range v.BlockedWords {
		if strings.Contains(strings.ToLower(cleaned), strings.ToLower(word)) {
			return nil, fmt.Errorf("company name contains blocked word: %s", word)
		}
	}

	// Validate confidence score
	if response.Company.Confidence < v.MinConfidence {
		return nil, fmt.Errorf("confidence score %.2f is below minimum %.2f",
			response.Company.Confidence, v.MinConfidence)
	}

	if response.Company.Confidence > v.MaxConfidence {
		return nil, fmt.Errorf("confidence score %.2f is above maximum %.2f",
			response.Company.Confidence, v.MaxConfidence)
	}

	// Return validated result
	return &CompanyIdentification{
		Name:       cleaned,
		Confidence: response.Company.Confidence,
	}, nil
}

func (v *CompanyResponseValidator) GetExpectedSchema() string {
	return `{
  "company": {
    "name": "Example Company",
    "confidence": 0.95
  }
}`
}
