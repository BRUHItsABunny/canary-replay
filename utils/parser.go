package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func ParseHCY(location string, omitHeaders map[string]struct{}) (*http.Request, error) {
	if !strings.HasSuffix(location, pathSeparator) {
		location += pathSeparator
	}
	reqSummary, err := openAndParseSummary(location + "request.json")
	if err != nil {
		return nil, fmt.Errorf("openAndParseSummary: %w", err)
	}
	var bodyBuf []byte
	reqContent, err := os.Open(location + "request_body.bin")
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}

	bodyBuf, err = io.ReadAll(reqContent)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	result, err := http.NewRequest(reqSummary.Method, reqSummary.URL, bytes.NewBuffer(bodyBuf))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	// Add headers but don't automatically add capital letters, API may expect lowercase only so maintain the actual captured case
	for headerName, headerValue := range reqSummary.Headers {
		_, omit := omitHeaders[strings.ToLower(headerName)]
		if !omit {
			result.Header[headerName] = []string{headerValue}
		}
	}
	return result, err
}

func openAndParseSummary(filePath string) (*CanaryJSON, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	result := new(CanaryJSON)
	err = json.Unmarshal(fileBytes, result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return result, err
}

func hasBody(method string) bool {
	methods := []string{
		"POST",
		"PATCH",
		"PUT",
	}
	for _, match := range methods {
		if match == method {
			return true
		}
	}
	return false
}
