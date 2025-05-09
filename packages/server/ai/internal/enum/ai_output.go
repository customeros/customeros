package enum

type AIOutputFormat string

const (
	AIOutputText AIOutputFormat = "text"
	AIOutputJson AIOutputFormat = "json"
)

func (a AIOutputFormat) String() string {
	return string(a)
}
