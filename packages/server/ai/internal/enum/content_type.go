package enum

type ContentType string

const (
	ContentArticle                ContentType = "article"
	ContentWhitepaper             ContentType = "whitepaper"
	ContentWebinar                ContentType = "webinar"
	ContentCaseStudy              ContentType = "case study"
	ContentProductPage            ContentType = "product page"
	ContentSolutionPage           ContentType = "solution page"
	ContentTestimonial            ContentType = "testimonial"
	ContentResearchReport         ContentType = "research report"
	ContentTechnicalDocumentation ContentType = "technical documentation"
	ContentUnknown                ContentType = ""
)
