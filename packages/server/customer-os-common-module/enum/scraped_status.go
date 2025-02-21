package enum

type ScrapeStatus string

const (
	ScrapeCompleted  ScrapeStatus = "COMPLETED"
	ScrapeError      ScrapeStatus = "ERROR"
	ScrapeNotScraped ScrapeStatus = "NOT_SCRAPED"
)

func (s ScrapeStatus) String() string {
	return string(s)
}
