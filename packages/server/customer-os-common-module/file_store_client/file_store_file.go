package fsc

type FileDTO struct {
	Id          string `json:"id"`
	FileName    string `json:"fileName"`
	MimeType    string `json:"mimeType"`
	Size        int64  `json:"size"`
	MetadataUrl string `json:"previewUrl"`
	DownloadUrl string `json:"downloadUrl"`
}
