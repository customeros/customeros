package model

import "time"

type Attachment struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	MimeType  string    `json:"mimeType"`
	FileName  string    `json:"fileName"`
	BasePath  string    `json:"basePath"`
	Size      int64     `json:"size"`
}

type AttachmentCreateResponse struct {
	Attachment `json:"attachment_Create"`
}

type AttachmentResponse struct {
	Attachment `json:"attachment"`
}
