package model

type CommentData struct {
	BaseData
	Content        string          `json:"content,omitempty"`
	ContentType    string          `json:"contentType,omitempty"`
	AuthorUser     ReferencedUser  `json:"authorUser,omitempty"`
	CommentedIssue ReferencedIssue `json:"commentedIssue,omitempty"`
}

func (l *CommentData) Normalize() {
	l.BaseData.Normalize()
	l.SetTimes()
}
