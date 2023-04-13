package model

import "time"

type Attachment struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	MimeType  string    `json:"mimeType"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	Extension string    `json:"extension"`
}

type AttachmentCreateResponse struct {
	Attachment `json:"attachment_Create"`
}

type AttachmentResponse struct {
	Attachment `json:"attachment"`
}
