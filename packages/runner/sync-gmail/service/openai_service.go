package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/repository"
	commonEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type EmailClassification struct {
	Email string

	IsPersonalEmail   bool
	PersonalFirstName string
	PersonalLastName  string

	IsOrganizationEmail bool
	OrganizationName    string
}

type openAiService struct {
	cfg          *config.Config
	repositories *repository.Repositories
}

type OpenAiService interface {
	FetchEmailsClassification(tenant, elementId, from string, to []string, cc []string, bcc []string) ([]*EmailClassification, error)
}

func (s *openAiService) FetchEmailsClassification(tenant, elementId, from string, to []string, cc []string, bcc []string) ([]*EmailClassification, error) {
	ctx := context.Background()

	var allEmails []string
	allEmails = append(allEmails, from)
	for _, prop := range [][]string{to, cc, bcc} {
		allEmails = append(allEmails, prop...)
	}

	var aiClassificationList []*EmailClassification
	for _, email := range allEmails {
		domainNode, err := s.repositories.DomainRepository.GetDomain(ctx, extractDomain(email))
		if err != nil {
			logrus.Errorf("failed to get domain for email %v :%v", email, err)
			return nil, err
		}
		if domainNode == nil {
			aiClassificationList = append(aiClassificationList, &EmailClassification{
				Email: email,
			})
		}
	}

	if aiClassificationList == nil || len(aiClassificationList) == 0 {
		return nil, nil
	}

	aiEmails := make([]string, len(aiClassificationList))
	for i, aiClassification := range aiClassificationList {
		aiEmails[i] = aiClassification.Email
	}
	emailsAsString := strings.Join(aiEmails, ", ")

	result, err := s.queryOpenAi(tenant, elementId, emailsAsString)
	if err != nil {
		logrus.Errorf("failed to query open ai: %v", err)
		return nil, err
	}

	choices := result["choices"].([]interface{})
	if len(choices) > 0 {

		if choices[0].(map[string]interface{})["finish_reason"] == "length" {
			logrus.Errorf("not enough token to generate ai classification: %v", err)
			return nil, err
		}

		categorization := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

		aiClasificationArray, err := ParseEmailClassifications(categorization)
		if err != nil {
			logrus.Errorf("failed to parse email classification: %v", err)
			return nil, err
		}

		for _, aiClasification := range aiClasificationArray {
			for _, e := range aiClassificationList {
				if e.Email == aiClasification.Email {
					e.IsPersonalEmail = aiClasification.IsPersonalEmail
					e.PersonalFirstName = aiClasification.PersonalFirstName
					e.PersonalLastName = aiClasification.PersonalLastName
					e.IsOrganizationEmail = aiClasification.IsOrganizationEmail
					e.OrganizationName = aiClasification.OrganizationName
				}
			}
		}
	}

	return aiClassificationList, nil
}

func (s *openAiService) queryOpenAi(tenant, elementId, emailsAsString string) (map[string]interface{}, error) {
	prompt := "For the emails in array: [" + emailsAsString + "] \n- i want to know if the email would be a personal address or a company address\n- if it's a company address, would it be a generic address like sales, no-reply, etc or a person in the company\n- format the response to be an array with JSON objects as the structure below with no other words added\n- you can use internet information\n- do not invent\n" +
		"\n{\n  \"Email\": \"FILL-THE-EMAIL-HERE\",\n  \"IsPersonalEmail\": false,\n  \"PersonalFirstName\": \"\",\n  \"PersonalLastName\": \"\",\n  \"IsOrganizationEmail\": false,\n  \"OrganizationName\": \"\"\n}"

	requestData := map[string]interface{}{}
	requestData["model"] = "gpt-3.5-turbo-16k"
	requestData["prompt"] = prompt
	requestData["maxTokensToSample"] = 1024

	requestBody, _ := json.Marshal(requestData)
	request, _ := http.NewRequest("POST", s.cfg.OpenAi.ApiPath+"/ask", strings.NewReader(string(requestBody)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Openline-API-KEY", s.cfg.OpenAi.ApiKey)

	nodeLabel := "RawEmail"
	log := commonEntity.AiPromptLog{
		CreatedAt:  time.Time{},
		AppSource:  "sync-gmail",
		Provider:   "anthropic",
		Model:      "claude-2",
		PromptType: "EMAIL_CLASSIFICATION",
		Tenant:     &tenant,
		NodeId:     &elementId,
		NodeLabel:  &nodeLabel,
		Prompt:     prompt,
	}
	storeLogId, errLog := s.repositories.CommonRepositories.AiPromptLogRepository.Store(log)
	if errLog != nil {
		fmt.Println("Error:", errLog)
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if response == nil || response.Body == nil {
		s.repositories.CommonRepositories.AiPromptLogRepository.UpdateError(storeLogId, "no response body")
		fmt.Println("Error making the API request: no response body")
		return nil, err
	}

	//store response
	bodyBytes, errBody := ioutil.ReadAll(response.Body)
	if errBody != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	bodyString := string(bodyBytes)
	err = s.repositories.CommonRepositories.AiPromptLogRepository.UpdateResponse(storeLogId, bodyString)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	if err != nil {
		s.repositories.CommonRepositories.AiPromptLogRepository.UpdateError(storeLogId, err.Error())
		fmt.Println("Error making the API request:", err)
		return nil, err
	}
	if response.StatusCode != 200 {
		s.repositories.CommonRepositories.AiPromptLogRepository.UpdateError(storeLogId, err.Error())
		fmt.Println("Error making the API request:", response.Status)
		return nil, fmt.Errorf("error making the API request: %s", response.Status)
	}
	defer response.Body.Close()

	var result map[string]interface{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		s.repositories.CommonRepositories.AiPromptLogRepository.UpdateError(storeLogId, err.Error())
		fmt.Println("Error parsing the API response:", err)
		return nil, err
	}
	return result, nil
}

func ParseEmailClassifications(input string) ([]EmailClassification, error) {
	var emails []EmailClassification

	// Try to unmarshal the input JSON string as an array of EmailClassification
	err := json.Unmarshal([]byte(input), &emails)
	if err == nil {
		return emails, nil
	}

	// If unmarshaling as an array failed, check if the input contains multiple JSON objects in a single string
	objects := strings.Split(input, "}\n{")
	if len(objects) > 1 {
		// Append the braces to each object and try to unmarshal each object as an EmailClassification
		for _, obj := range objects {
			if !strings.HasPrefix(obj, "{") {
				obj = "{" + obj
			}
			if !strings.HasSuffix(obj, "}") {
				obj = obj + "}"
			}
			var email EmailClassification
			err := json.Unmarshal([]byte(obj), &email)
			if err != nil {
				return nil, fmt.Errorf("failed to parse JSON: %v", err)
			}
			emails = append(emails, email)
		}
		return emails, nil
	}

	// If it's neither an array nor multiple JSON objects in a single string, try to unmarshal as a single EmailClassification
	var email EmailClassification
	err = json.Unmarshal([]byte(input), &email)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// If unmarshaling as a single EmailClassification succeeded, return it as a slice of EmailClassification
	return []EmailClassification{email}, nil
}

func NewOpenAiService(cfg *config.Config, repositories *repository.Repositories) OpenAiService {
	return &openAiService{
		cfg:          cfg,
		repositories: repositories,
	}
}
