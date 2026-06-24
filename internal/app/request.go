package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    map[string]any    `json:"body"`
}

func (r *Request) Send(ctx *TestContext, baseURL string) (*http.Response, time.Duration, error) {
	err := FindAndReplaceVars(ctx, r)
	if err != nil {
		return nil, 0, err
	}

	url := baseURL + r.URL

	req, err := http.NewRequest(strings.ToUpper(r.Method), url, r.normilizedBody())

	if err != nil {
		fmt.Printf("Error create request %v\n", err.Error())
		return nil, 0, err
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}

	start := time.Now()

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error exec request %v\n", err.Error())
		return nil, 0, err
	}

	return resp, time.Since(start), nil
}

func (r *Request) normilizedBody() *bytes.Reader {
	jsonBody, _ := json.Marshal(r.Body)
	return bytes.NewReader(jsonBody)
}
