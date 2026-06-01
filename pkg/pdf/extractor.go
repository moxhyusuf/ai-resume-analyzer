package pdf

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

// ExtractText reads a PDF from raw bytes and returns extracted plain text.
func ExtractText(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	pdfReader, err := pdf.NewReader(reader, int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("open pdf: %w", err)
	}

	var sb strings.Builder
	for i := 1; i <= pdfReader.NumPage(); i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			continue // skip unreadable pages
		}
		sb.WriteString(text)
		sb.WriteString("\n")
	}

	result := strings.TrimSpace(sb.String())
	if result == "" {
		return "", fmt.Errorf("no readable text found in PDF")
	}

	// Limit to ~8000 chars to stay within token limits
	if len(result) > 8000 {
		result = result[:8000]
	}

	return result, nil
}
