package plugin

import (
	"bytes"
	"net/http"
	"strings"
)

func SplitBytes(data []byte) [][]byte {
	var result [][]byte
	interval := 65535

	i := 0
	for i <= len(data) {
		end := i + interval
		// Make sure we don't go beyond the array's length
		if end > len(data) {
			end = len(data)
		}
		// Append the chunk to the result
		result = append(result, data[i:end])
		i += interval
	}

	return result
}

func InvertSplitBytes(chunks [][]byte) []byte {
	var buffer bytes.Buffer
	for _, chunk := range chunks {
		buffer.Write(chunk)
	}
	return buffer.Bytes()
}

func DetermineMimeType(fileName string, content []byte) string {
	if strings.Contains(fileName, ".docx") {
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	}
	if strings.Contains(fileName, ".xlsx") {
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}
	return http.DetectContentType(content)
}
