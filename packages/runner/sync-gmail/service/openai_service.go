package service

import (
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

type openAiService struct {
	cfg          *config.Config
	repositories *repository.Repositories
}

type OpenAiService interface {
	AskForOrganizationNameByDomain(tenant, elementId, domain string) (string, error)
}

func (s *openAiService) AskForOrganizationNameByDomain(tenant, elementId, domain string) (string, error) {
	result, err := s.queryOpenAi(tenant, elementId, domain)
	if err != nil {
		logrus.Errorf("failed to query open ai: %v", err)
		return "", err
	}

	choices := result["choices"].([]interface{})
	if len(choices) > 0 {

		if choices[0].(map[string]interface{})["finish_reason"] == "length" {
			logrus.Errorf("not enough token to generate ai classification: %v", err)
			return "", err
		}

		return choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string), nil
	}

	return "", nil
}

func (s *openAiService) queryOpenAi(tenant, elementId, domain string) (map[string]interface{}, error) {
	prompt := "What name would have the organization that has this domain: " + domain + "\nI want a simple response in a single line with no other words ar analysis."

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
		PromptType: "ORGANIZATION_NAME_BY_DOMAIN",
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

func NewOpenAiService(cfg *config.Config, repositories *repository.Repositories) OpenAiService {
	return &openAiService{
		cfg:          cfg,
		repositories: repositories,
	}
}
