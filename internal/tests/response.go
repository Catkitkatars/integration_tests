package tests

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Response struct {
	Body         string        `json:"body"`
	StatusCode   int           `json:"statusCode"`
	Headers      http.Header   `json:"headers"`
	ResponseTime time.Duration `json:"responseTime"`
}

func NewResponse(resp *http.Response, respTime time.Duration) (*Response, error) {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Faild to read response body")
	}

	return &Response{
		Body:         string(bodyBytes),
		StatusCode:   resp.StatusCode,
		Headers:      resp.Header,
		ResponseTime: respTime,
	}, nil
}
