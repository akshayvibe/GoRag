package handler

import (
	"encoding/json"

	"github.com/akshayvibe/GoRag/internal/indexing/upload"
	// "github.com/akshayvibe/GoRag/internal/models"
)


func UploadHandler(path string) (string, error) {
	pdfContent, err := upload.OpenPdf(path)
	if err != nil {
		return "", err
	}
	data,err:=json.MarshalIndent(pdfContent,"","");
	if err != nil {
		return "", err
	}
	return string(data), nil
}