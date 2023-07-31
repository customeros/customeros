package organization

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type DomainScraper struct {
	log logger.Logger
	cfg *config.Config
}

func NewDomainScraper(log logger.Logger, cfg *config.Config) *DomainScraper {
	return &DomainScraper{
		log: log,
		cfg: cfg,
	}
}

func (ds *DomainScraper) Scrape(domain string) (*WebscrapeResponseV1, error) {
	domainUrl := "https://" + domain
	jsonStruct := jsonStructure()

	html, err := ds.getHtml(domainUrl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to getHtml domain")
	}
	socialLinks, err := ds.extractSocialLinks(html)
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract social links")
	}

	text, err := ds.extractRelevantText(html)
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract relevant text")
	}

	companyAnalysis, err := ds.runCompanyPrompt(text)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run company prompt")
	}

	r, err := ds.runDataPrompt(companyAnalysis, &domainUrl, socialLinks, jsonStruct)

	if err != nil {
		return nil, errors.Wrap(err, "failed to run data prompt")
	}
	return r, nil
}

func (ds *DomainScraper) getHtml(domainUrl string) (*string, error) {
	response, err := ds.getRequest(ds.cfg.Services.ScrapingBeeApiKey, domainUrl)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get response")
	}

	// Read Response Body
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	s := string(respBody)

	return &s, nil
}

func (ds *DomainScraper) runCompanyPrompt(text *string) (*string, error) {
	prompt := strings.ReplaceAll(ds.cfg.Services.OpenAi.ScrapeCompanyPrompt, "{{text}}", *text)
	ds.log.Printf("prompt: %s", prompt)

	aiResult, err := ds.openai(prompt)

	if err != nil {
		return nil, errors.Wrap(err, "unable to get openai result")
	}

	return aiResult, nil
}

func (ds *DomainScraper) getRequest(api_key string, domainUrl string) (*http.Response, error) {
	// Create client
	client := &http.Client{}

	urlEscaped := url.QueryEscape(domainUrl) // Encoding the URL
	// Create request
	req, err := http.NewRequest("GET", "https://app.scrapingbee.com/api/v1/?api_key="+api_key+"&url="+urlEscaped, nil) // Create the request the request
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		return nil, errors.Wrap(parseFormErr, "failed to parse form")
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch request")
	}
	return resp, nil // Return the response
}

func (ds *DomainScraper) extractRelevantText(html *string) (*string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create document from reader")
	}

	doc.Find("script, style").Remove()

	var texts []string
	doc.Find("*").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return s.Children().Length() == 0
	}).Each(func(_ int, s *goquery.Selection) {
		text := s.Text()
		text = strings.TrimSpace(text)
		if text != "" && !contains(texts, text) {
			texts = append(texts, text)
		}
	})

	text := strings.Join(texts, " ")
	ds.log.Printf("text: %s", text)
	return &text, nil
}

func contains(slice []string, text string) bool {
	for _, s := range slice {
		if s == text {
			return true
		}
	}
	return false
}

func (ds *DomainScraper) extractSocialLinks(html *string) (*string, error) {
	socialSites := map[string]string{
		"linkedin":  "linkedin.com",
		"twitter":   "twitter.com",
		"instagram": "instagram.com",
		"facebook":  "facebook.com",
		"youtube":   "youtube.com",
		"github":    "github.com",
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create document from reader")
	}
	links := make(map[string]string)

	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		for name, site := range socialSites {
			if strings.Contains(href, site) {
				links[name] = href
				break
			}
		}
	})

	socialLinks := fmt.Sprintf("%s", links)
	ds.log.Printf("social links: %s", socialLinks)
	return &socialLinks, nil
}

func (ds *DomainScraper) openai(prompt string) (*string, error) {
	requestData := map[string]interface{}{}
	requestData["model"] = "TextDavinci003"
	requestData["prompt"] = prompt

	requestBody, _ := json.Marshal(requestData)
	request, _ := http.NewRequest("POST", ds.cfg.Services.OpenAi.ApiPath+"/ask", strings.NewReader(string(requestBody)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Openline-API-KEY", ds.cfg.Services.OpenAi.ApiKey)

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
		return &categorization, nil
	}
	return nil, errors.New("no result found")
}

func (ds *DomainScraper) runDataPrompt(analysis, domainUrl, socials, jsonStructure *string) (*WebscrapeResponseV1, error) {

	replacements := map[string]string{
		"{{ANALYSIS}}":       *analysis,
		"{{DOMAIN_URL}}":     *domainUrl,
		"{{SOCIALS}}":        *socials,
		"{{JSON_STRUCTURE}}": *jsonStructure,
	}

	sPrompt := ds.cfg.Services.OpenAi.ScrapeDataPrompt
	for k, v := range replacements {
		sPrompt = strings.ReplaceAll(sPrompt, k, v)
	}

	cleaned, err := ds.openai(sPrompt)
	if err != nil {
		return nil, err
	}
	scrapeResponse := WebscrapeResponseV1{}
	err = json.Unmarshal([]byte(*cleaned), &scrapeResponse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}

	return &scrapeResponse, nil
}

func jsonStructure() *string {
	data := WebscrapeResponseV1{
		CompanyName:      "...",
		Website:          "...",
		Market:           "...",
		Industry:         "...",
		IndustryGroup:    "...",
		SubIndustry:      "...",
		TargetAudience:   "...",
		ValueProposition: "...",
		Github:           "...",
		Linkedin:         "...",
		Twitter:          "...",
		Youtube:          "...",
		Instagram:        "...",
		Facebook:         "...",
	}

	jsonStructure, _ := json.Marshal(data)

	s := string(jsonStructure)
	return &s
}
