package model

import "time"

type VConAppended struct {
	UUID string `json:"uuid"`
}

type VConParty struct {
	Tel       *string `json:"tel,omitempty"`
	Stir      *string `json:"stir,omitempty"`
	Mailto    *string `json:"mailto,omitempty"`
	Name      *string `json:"name,omitempty"`
	ContactId *string `json:"contact_id,omitempty"`
	UserId    *string `json:"user_id,omitempty"`
}

type VConAnalysisType string

const (
	SUMMARY     VConAnalysisType = "summary"
	TRANSCRIPT  VConAnalysisType = "transcript"
	TRANSLATION VConAnalysisType = "translation"
	SENTIMENT   VConAnalysisType = "sentiment"
	TTS         VConAnalysisType = "tts"
)

type VConAnalysis struct {
	Type     VConAnalysisType `json:"type"`
	Dialog   []int64          `json:"dialog"`
	MimeType string           `json:"mimetype"`
	Body     string           `json:"body"`
	Encoding VConEncoding     `json:"encoding"`
}

type VConDialogType string

const (
	TEXT      VConDialogType = "text"
	RECORDING VConDialogType = "recording"
)

type VConEncoding string

const (
	NONE   VConEncoding = "None"
	BASE64 VConEncoding = "base64url"
	JSON   VConEncoding = "json"
)

type VConDialog struct {
	Type     VConDialogType `json:"type"`
	Start    time.Time      `json:"start"`
	Duration int64          `json:"duration"`
	Parties  []int64        `json:"parties"`
	MimeType string         `json:"mimetype"`
	Body     string         `json:"body"`
	Encoding VConEncoding   `json:"encoding"`
}

type VConAttachment struct {
	Parties  []int64      `json:"parties"`
	MimeType string       `json:"mimetype"`
	Body     string       `json:"body"`
	Encoding VConEncoding `json:"encoding"`
}

type VCon struct {
	VCon        string           `json:"vcon,default=0.0.1"`
	UUID        string           `json:"uuid"`
	Subject     string           `json:"subject"`
	Appended    *VConAppended    `json:"appended,omitempty"`
	Parties     []VConParty      `json:"parties,omitempty"`
	Dialog      []VConDialog     `json:"dialog,omitempty"`
	Attachments []VConAttachment `json:"attachments,omitempty"`
	Analysis    []VConAnalysis   `json:"analysis,omitempty"`
}
