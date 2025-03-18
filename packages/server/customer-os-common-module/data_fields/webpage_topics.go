package data_fields

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Topic struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

type WebpageTopicsWithConfidence struct {
	Topics []Topic `json:"topics"`
}

type TopicsResponseValidator struct {
	MaxTopics     int      // Maximum number of topics allowed
	MinTopics     int      // Minimum number of topics required
	MaxTopicLen   int      // Maximum length of each topic string
	DisallowEmpty bool     // Whether to disallow empty topics
	BlockedWords  []string // Words that shouldn't appear in topics
	MinConfidence float64  // Minimum confidence score allowed
	MaxConfidence float64  // Maximum confidence score allowed
}

func NewTopicsResponseValidator() *TopicsResponseValidator {
	return &TopicsResponseValidator{
		MaxTopics:     5,
		MinTopics:     1,
		MaxTopicLen:   100,
		DisallowEmpty: true,
		BlockedWords:  []string{},
		MinConfidence: 0.0,
		MaxConfidence: 1.0,
	}
}

func (v *TopicsResponseValidator) IsValidJSON(jsonStr string) bool {
	var topics WebpageTopicsWithConfidence
	err := json.Unmarshal([]byte(jsonStr), &topics)
	return err == nil && topics.Topics != nil
}

func (v *TopicsResponseValidator) ValidateResponse(jsonStr string) ([]Topic, error) {
	// Parse the JSON
	var response WebpageTopicsWithConfidence
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	// Check if topics exist
	if response.Topics == nil {
		return nil, fmt.Errorf("missing 'topics' array in response")
	}

	// Check number of topics
	if len(response.Topics) > v.MaxTopics {
		return nil, fmt.Errorf("too many topics: got %d, maximum allowed is %d",
			len(response.Topics), v.MaxTopics)
	}

	if len(response.Topics) < v.MinTopics {
		return nil, fmt.Errorf("too few topics: got %d, minimum required is %d",
			len(response.Topics), v.MinTopics)
	}

	// Validate each topic
	var cleanedTopics []Topic
	seenTopics := make(map[string]bool)

	for i, topic := range response.Topics {
		// Clean the topic name
		cleaned := strings.TrimSpace(topic.Name)

		// Check if empty
		if v.DisallowEmpty && cleaned == "" {
			return nil, fmt.Errorf("empty topic name at index %d", i)
		}

		// Check length
		if len(cleaned) > v.MaxTopicLen {
			return nil, fmt.Errorf("topic name at index %d exceeds maximum length: got %d, max is %d",
				i, len(cleaned), v.MaxTopicLen)
		}

		// Check for blocked words
		for _, word := range v.BlockedWords {
			if strings.Contains(strings.ToLower(cleaned), strings.ToLower(word)) {
				return nil, fmt.Errorf("topic at index %d contains blocked word: %s", i, word)
			}
		}

		// Validate confidence score
		if topic.Confidence < v.MinConfidence {
			return nil, fmt.Errorf("topic '%s' has confidence score %.2f below minimum %.2f",
				cleaned, topic.Confidence, v.MinConfidence)
		}

		if topic.Confidence > v.MaxConfidence {
			return nil, fmt.Errorf("topic '%s' has confidence score %.2f above maximum %.2f",
				cleaned, topic.Confidence, v.MaxConfidence)
		}

		// Check for duplicates
		topicLower := strings.ToLower(cleaned)
		if seenTopics[topicLower] {
			return nil, fmt.Errorf("duplicate topic found: %s", cleaned)
		}
		seenTopics[topicLower] = true

		// Add to cleaned topics
		cleanedTopics = append(cleanedTopics, Topic{
			Name:       cleaned,
			Confidence: topic.Confidence,
		})
	}

	return cleanedTopics, nil
}

func (v *TopicsResponseValidator) GetExpectedSchema() string {
	return `{
  "topics": [
    {
      "name": "topic 1",
      "confidence": 0.95
    },
    {
      "name": "topic 2",
      "confidence": 0.87
    },
    {
      "name": "topic 3",
      "confidence": 0.72
    }
  ]
}`
}
