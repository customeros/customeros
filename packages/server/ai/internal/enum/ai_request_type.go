package enum

type AIRequestType string

const (
	AIRequestGeneric            AIRequestType = "generic"
	AIRequestCompanyDescription AIRequestType = "company_description"
	AIRequestCompanyName        AIRequestType = "company_name"
	AIRequestContentStage       AIRequestType = "content_journey_stage"
	AIRequestEmail              AIRequestType = "email"
	AIRequestIndustryCode       AIRequestType = "industry_code"
	AIRequestWebpageCategory    AIRequestType = "webpage_category"
	AIRequestWebpageTopics      AIRequestType = "webpage_topics"
)

func (a AIRequestType) String() string {
	return string(a)
}
