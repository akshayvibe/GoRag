package upload

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/akshayvibe/GoRag/internal/models"
	"github.com/dslipak/pdf"
	"github.com/google/uuid"
)

type PdfContent struct {
	Text     string
	Metadata models.Metadata
}

func OpenPdf(path string) (*PdfContent, error) {
	// Get PDF content
	content, err := readPdf(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read pdf: %w", err)
	}

	// Get document metadata
	metadata := getMetadata(path)

	return &PdfContent{
		Text:     content,
		Metadata: metadata,
	}, nil
}
func readPdf(path string) (string, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("open pdf: %w", err)
	}
	defer r.Trailer().Reader().Close()

	b, err := r.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("extract pdf text: %w", err)
	}

	var buf bytes.Buffer

	_, err = buf.ReadFrom(b)
	if err != nil {
		return "", fmt.Errorf("read pdf content: %w", err)
	}

	return buf.String(), nil
}

func getMetadata(path string) models.Metadata {
	return models.Metadata{
		DocumentID: uuid.New().String(),
		FileName:   filepath.Base(path),
	}
}