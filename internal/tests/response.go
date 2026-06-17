package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Response struct {
	Body         map[string]any `json:"body"`
	StatusCode   int            `json:"statusCode"`
	Headers      http.Header    `json:"headers"`
	ResponseTime time.Duration  `json:"responseTime"`
}

func NewResponse(resp *http.Response, respTime time.Duration) (*Response, error) {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Faild to read response body")
	}
	var body map[string]any

	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			return nil, fmt.Errorf("Cannot unmarshal json body")
		}
	}

	return &Response{
		Body:         body,
		StatusCode:   resp.StatusCode,
		Headers:      resp.Header,
		ResponseTime: respTime,
	}, nil
}
