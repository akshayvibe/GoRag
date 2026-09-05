package upload

import (
	"bytes"

	"github.com/dslipak/pdf"
)

func OpenPdf(path string) (string, error) {
	content, err := readPdf(path)
	if err != nil {
		return "", err
	}

	return content, nil
}

func readPdf(path string) (string, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}

	defer r.Trailer().Reader().Close()

	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	_, err = buf.ReadFrom(b)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}