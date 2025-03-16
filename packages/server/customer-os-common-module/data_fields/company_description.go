package data_fields

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

type CompanyDescription struct {
	Description string `json:"description"`
}

type CompanyDescriptionResponse struct {
	Company CompanyDescription `json:"company"`
}

type CompanyDescriptionResponseValidator struct {
	MaxLength     int      // Maximum length in characters
	MinLength     int      // Minimum length in characters
	DisallowEmpty bool     // Whether to disallow empty description
	BlockedWords  []string // Words that shouldn't appear in description
	RequireWords  []string // Words that must appear in description
}

func NewCompanyDescriptionResponseValidator() *CompanyDescriptionResponseValidator {
	return &CompanyDescriptionResponseValidator{
		MaxLength:     400, // Allow up to 400 chars, though we ask AI for 300 to be conservative
		MinLength:     50,  // Reasonable minimum for a useful description
		DisallowEmpty: true,
		BlockedWords:  []string{},
		RequireWords:  []string{},
	}
}

func (v *CompanyDescriptionResponseValidator) IsValidJSON(jsonStr string) bool {
	var response CompanyDescriptionResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	return err == nil
}

func (v *CompanyDescriptionResponseValidator) ValidateResponse(jsonStr string) (*CompanyDescription, error) {
	// Parse the JSON
	var response CompanyDescriptionResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	// Clean the description
	cleaned := strings.TrimSpace(response.Company.Description)

	// Check if empty
	if v.DisallowEmpty && cleaned == "" {
		return nil, fmt.Errorf("empty company description")
	}

	// Check description length
	charCount := utf8.RuneCountInString(cleaned)
	if charCount > v.MaxLength {
		return nil, fmt.Errorf("description exceeds maximum length: got %d characters, max is %d",
			charCount, v.MaxLength)
	}

	if charCount < v.MinLength {
		return nil, fmt.Errorf("description is too short: got %d characters, min is %d",
			charCount, v.MinLength)
	}

	// Check for blocked words
	for _, word := range v.BlockedWords {
		if strings.Contains(strings.ToLower(cleaned), strings.ToLower(word)) {
			return nil, fmt.Errorf("description contains blocked word: %s", word)
		}
	}

	// Check for required words
	for _, word := range v.RequireWords {
		if !strings.Contains(strings.ToLower(cleaned), strings.ToLower(word)) {
			return nil, fmt.Errorf("description does not contain required word: %s", word)
		}
	}

	// Verify it's a full sentence with proper punctuation
	if !strings.HasSuffix(cleaned, ".") && !strings.HasSuffix(cleaned, "!") && !strings.HasSuffix(cleaned, "?") {
		return nil, fmt.Errorf("description should end with proper punctuation")
	}

	// Return validated result
	return &CompanyDescription{
		Description: cleaned,
	}, nil
}

func (v *CompanyDescriptionResponseValidator) GetExpectedSchema() string {
	return `{
  "company": {
    "description": "Example Company provides B2B sales automation solutions to small businesses and help streamline their lead generation and customer engagement processes. Their CRM platform uses AI to qualify leads, automate outreach, and provide analytics to improve conversion rates."
  }
}`
}
