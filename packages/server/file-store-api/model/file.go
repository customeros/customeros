package model

type File struct {
	ID       string
	FileName string
	MimeType string
	BasePath string
	Size     int64
	CdnUrl   string
}
