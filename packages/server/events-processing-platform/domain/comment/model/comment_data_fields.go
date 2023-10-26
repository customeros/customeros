package model

type CommentDataFields struct {
	Content          string
	ContentType      string
	AuthorUserId     *string
	CommentedIssueId *string
}
