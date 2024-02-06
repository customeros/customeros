package utils

import (
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"mime/multipart"
)

func GetFileTypeHeadFromBytes(bytes *[]byte) ([]byte, error) {
	head := make([]byte, 1024)
	copy(head, *bytes)
	return head, nil
}

func GetFileTypeHeadFromMultipart(file multipart.File) ([]byte, error) {
	head := make([]byte, 1024)
	_, err := file.Read(head)
	if err != nil {
		return nil, err
	}
	return head, nil
}

func GetFileType(head []byte) (types.Type, error) {
	//TODO docx is not recognized
	//https://github.com/h2non/filetype/issues/121
	kind, _ := filetype.Match(head)

	return kind, nil
}
