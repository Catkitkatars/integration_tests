package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"integrationstests/internal/app/input"
	"io"
	"net/http"
	"regexp"
	"strings"
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

type Request struct {
	method  string
	url     string
	headers map[string]string
	body    map[string]any
}

func NewRequest(
	method string,
	url string,
	headers map[string]string,
	body map[string]any,
) *Request {
	return &Request{
		method:  method,
		url:     url,
		headers: headers,
		body:    body,
	}
}

func InitRequest(ctx context.Context, rq input.RequestData) (*Request, error) {
	if rq.Method == "" {
		return nil, fmt.Errorf("request method is required")
	}

	method := strings.ToUpper(rq.Method)
	if !isValidMethod(method) {
		return nil, fmt.Errorf("unsupported request method: %s", rq.Method)
	}

	return NewRequest(
		method,
		rq.URL,
		rq.Headers,
		rq.Body,
	), nil
}

func (r *Request) Send() (*Response, error) {
	req, err := http.NewRequest(strings.ToUpper(r.method), r.url, r.normilizedBody())

	if err != nil {
		fmt.Printf("Error create request %v\n", err.Error())
		return nil, err
	}

	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}

	start := time.Now()

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error exec request %v\n", err.Error())
		return nil, err
	}

	response, err := NewResponse(resp, time.Since(start))

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *Request) normilizedBody() *bytes.Reader {
	jsonBody, _ := json.Marshal(r.body)
	return bytes.NewReader(jsonBody)
}

func isValidMethod(m string) bool {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	for _, v := range methods {
		if m == v {
			return true
		}
	}

	return false
}

func (r *Request) FindAndReplaceVars(ctx context.Context) error {
	url, err := r.replaceVars(ctx, r.url, r.findVars(r.url))
	if err != nil {
		return err
	}
	r.url = url

	headers := make(map[string]string, len(r.headers))
	for key, value := range r.headers {
		header, err := r.replaceVars(ctx, value, r.findVars(value))
		if err != nil {
			return err
		}

		headers[key] = header
	}
	r.headers = headers

	body, err := r.replaceVarsFromBody(ctx, r.body)
	if err != nil {
		return err
	}

	r.body = body

	return nil
}

func (r *Request) replaceVarsFromBody(ctx context.Context, body map[string]any) (map[string]any, error) {
	handledBody := make(map[string]any, len(body))

	for key, value := range body {
		handledValue, err := r.replaceVarsFromValue(ctx, value)
		if err != nil {
			return nil, err
		}

		handledBody[key] = handledValue
	}

	return handledBody, nil
}

func (r *Request) replaceVarsFromValue(ctx context.Context, value any) (any, error) {
	switch value := value.(type) {
	case string:
		return r.replaceVars(ctx, value, r.findVars(value))

	case map[string]any:
		return r.replaceVarsFromBody(ctx, value)

	case []any:
		handledSlice := make([]any, 0, len(value))

		for _, item := range value {
			handledItem, err := r.replaceVarsFromValue(ctx, item)
			if err != nil {
				return nil, err
			}

			handledSlice = append(handledSlice, handledItem)
		}

		return handledSlice, nil

	default:
		return value, nil
	}
}

func (r *Request) findVars(text string) []string {
	matches := regexp.
		MustCompile(`\{\{([a-zA-Z0-9_]+)\}\}`).
		FindAllStringSubmatch(text, -1)

	var vars []string
	for _, match := range matches {
		vars = append(vars, match[1])
	}

	return vars
}

func (r *Request) replaceVars(ctx context.Context, text string, vars []string) (string, error) {
	result := text

	for _, name := range vars {
		value := ctx.Value(name)
		if value == nil {
			return result, fmt.Errorf("replace error: variable %q not found in context", name)
		}

		placeholder := "{{" + name + "}}"
		result = strings.ReplaceAll(result, placeholder, fmt.Sprint(value))
	}

	return result, nil
}
