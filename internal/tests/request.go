package tests

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
	Body    any               `json:"body"`
}

func (r *Request) Send(ctx *TestContext, url string) (*http.Response, time.Duration, error) {
	req, err := http.NewRequest(strings.ToUpper(r.Method), url, r.normilizedBody())

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	if err != nil {
		fmt.Printf("Error create request %v\n", err.Error())
		return nil, 0, err
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
