package interfaces

import "context"

type PostmarkService interface {
	SendNotification(ctx context.Context, postmarkEmail PostmarkEmail, tenant string) error
	CreateServerIfNotExists(ctx context.Context) error
	DeleteServer(ctx context.Context, tenant string) error
}

type PostmarkEmail struct {
	WorkflowId    string            `json:"workflowId"`
	MessageStream string            `json:"messageStream"`
	TemplateData  map[string]string `json:"templateData"`
	From          string            `json:"from"`
	To            string            `json:"to"`
	CC            []string          `json:"cc"`
	BCC           []string          `json:"bcc"`
	Subject       string            `json:"subject"`
	Attachments   []PostmarkEmailAttachment
}

type PostmarkEmailAttachment struct {
	Filename       string
	ContentEncoded string
	ContentType    string
	ContentID      string
}
