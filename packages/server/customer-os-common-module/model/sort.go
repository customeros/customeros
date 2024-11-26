package model

type SortBy struct {
	By            string           `json:"by"`
	Direction     SortingDirection `json:"direction"`
	CaseSensitive *bool            `json:"caseSensitive,omitempty"`
}

type SortingDirection string

const (
	SortingDirectionAsc  SortingDirection = "ASC"
	SortingDirectionDesc SortingDirection = "DESC"
)
