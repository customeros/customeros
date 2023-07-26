package service

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"net/http"
	"strings"
)

type OpenAiEmailClassification struct {
	Email string

	IsPersonalEmail   bool
	PersonalFirstName string
	PersonalLastName  string

	IsOrganizationEmail bool
	OrganizationName    string
}

type openAiService struct {
	cfg *config.Config
}

type OpenAiService interface {
	FetchEmailsClassification(from string, to []string, cc []string, bcc []string) ([]*OpenAiEmailClassification, error)
}

func (s *openAiService) FetchEmailsClassification(from string, to []string, cc []string, bcc []string) ([]*OpenAiEmailClassification, error) {
	classificationList := buildEmailsClassificationList(from, to, cc, bcc)

	// Send the request to OpenAI API
	categorizations := make(map[string]string)

	for _, e := range classificationList {
		email := strings.TrimSpace(e.Email)

		prompt := "For the email: " + email + " I want:\n- i want to know if the email would be a personal address or a company address\n- if it's a company address, would be generic address like sales, no-reply, etc or a person in the company\n- the response should be as below, no other words added\n- use internet information\n- do not invent \n" +
			"{\"Email\": \"" + email + "\"," +
			"\"IsPersonalEmail\": false," +
			"\"PersonalFirstName\": \"\"" +
			"\"PersonalLastName\": \"\"" +
			"\"IsOrganizationEmail\": false," +
			"\"OrganizationName\": \"\"" +
			"}"

		requestData := map[string]interface{}{}
		requestData["model"] = "gpt-4"
		requestData["prompt"] = prompt

		requestBody, _ := json.Marshal(requestData)
		request, _ := http.NewRequest("POST", s.cfg.OpenAi.ApiPath+"/ask", strings.NewReader(string(requestBody)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-Openline-API-KEY", s.cfg.OpenAi.ApiKey)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Error making the API request:", err)
			return nil, err
		}
		if response.StatusCode != 200 {
			fmt.Println("Error making the API request:", response.Status)
			return nil, fmt.Errorf("error making the API request: %s", response.Status)
		}
		defer response.Body.Close()

		var result map[string]interface{}
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			fmt.Println("Error parsing the API response:", err)
			return nil, err
		}

		choices := result["choices"].([]interface{})
		if len(choices) > 0 {
			categorization := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
			categorizations[email] = categorization
		}
	}

	for email, categorization := range categorizations {

		var m OpenAiEmailClassification

		err := json.Unmarshal([]byte(categorization), &m)
		if err != nil {
			fmt.Println("Error parsing the API response:", err)
			return nil, err
		}
		for _, e := range classificationList {
			if e.Email == email {
				e.IsPersonalEmail = m.IsPersonalEmail
				e.PersonalFirstName = m.PersonalFirstName
				e.PersonalLastName = m.PersonalLastName
				e.IsOrganizationEmail = m.IsOrganizationEmail
				e.OrganizationName = m.OrganizationName
			}
		}
	}

	return classificationList, nil
}

func buildEmailsClassificationList(from string, to []string, cc []string, bcc []string) []*OpenAiEmailClassification {
	var emailsClassification []*OpenAiEmailClassification
	emailsClassification = append(emailsClassification, &OpenAiEmailClassification{
		Email: from,
	})
	for _, email := range to {
		emailsClassification = append(emailsClassification, &OpenAiEmailClassification{
			Email: email,
		})
	}
	for _, email := range cc {
		emailsClassification = append(emailsClassification, &OpenAiEmailClassification{
			Email: email,
		})
	}
	for _, email := range bcc {
		emailsClassification = append(emailsClassification, &OpenAiEmailClassification{
			Email: email,
		})
	}
	return emailsClassification
}

func NewOpenAiService(cfg *config.Config) OpenAiService {
	return &openAiService{
		cfg: cfg,
	}
}
