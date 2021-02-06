package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func ParseHCY(location string, isHttps bool) (*http.Request, error) {
	var (
		reqContent *os.File
		reqSummary *CanaryJSON
		err        error
		result     *http.Request
	)

	if !strings.HasSuffix(location, pathSeparator) {
		location += pathSeparator
	}

	reqSummary, err = openAndParseSummary(location + "request.json")
	if err == nil {
		// TODO: this line might actually return nil?, may not need an if statement here
		reqContent, err = os.Open(location + "request_body.bin")
		if err == nil {
			result, err = http.NewRequest(reqSummary.Method, reqSummary.URL, reqContent)
		} else {
			result, err = http.NewRequest(reqSummary.Method, reqSummary.URL, nil)
		}

		if err == nil {
			// Add headers but don't automatically add capital letters, API may expect lowercase only so maintain the actual captured case
			headers := result.Header
			for headerName, headerValue := range reqSummary.Headers {
				headers[headerName] = []string{headerValue}
			}
			result.Header = headers
		}
	}

	return result, err
}

func openAndParseSummary(filePath string) (*CanaryJSON, error) {
	var (
		file      *os.File
		fileBytes []byte
		err       error
		result    = new(CanaryJSON)
	)

	file, err = os.Open(filePath)
	if err == nil {
		fileBytes, err = ioutil.ReadAll(file)
		if err == nil {
			err = json.Unmarshal(fileBytes, result)
		}
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
