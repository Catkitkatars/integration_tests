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

	"github.com/tidwall/gjson"
)

type StepInterface interface {
	Do(context.Context) (*StepResult, error)
}

type Step struct {
	name        string
	method      string
	url         string
	headers     string
	body        string
	extractData map[string]any
	asserts     *Asserts
}

type StepResult struct {
	AssertResults []*AssertsResult
}

func NewStep(
	name,
	method,
	url,
	headers,
	body string,
	asserts *Asserts,
	extract map[string]any) (*Step, error) {
	return &Step{
		name:        name,
		method:      method,
		url:         url,
		headers:     headers,
		body:        body,
		asserts:     asserts,
		extractData: extract,
	}, nil
}

func InitStep(data *input.StepData) (*Step, error) {
	if data.Request == nil {
		return nil, fmt.Errorf("unknown step type")
	}

	asserts := NewAsserts(data.Asserts.StatusCode, data.Asserts.Equals, data.Asserts.Exists)

	return NewStep(
		data.Name,
		data.Request.Method,
		data.Request.URL,
		data.Request.Headers,
		data.Request.Body,
		asserts,
		data.Extract)
}

func (s *Step) Do(ctx context.Context) (*StepResult, error) {
	err := s.prepearRequestData(ctx)
	if err != nil {
		return nil, fmt.Errorf("Step.Do(): %w", err)
	}

	resp, err := s.request()
	if err != nil {
		return nil, fmt.Errorf("Step.Do(): %w", err)
	}

	assertsRes, err := s.asserts.Check(resp)

	if err != nil {
		return nil, fmt.Errorf("Step.Do(): %w", err)
	}

	vars, err := s.extract(resp)

	if err != nil {
		return nil, fmt.Errorf("Step.Do(): %w", err)
	}

	for k, v := range vars {
		ctx = context.WithValue(ctx, k, v)
	}

	return &StepResult{
		AssertResults: assertsRes,
	}, nil
}

func (s *Step) request() (*http.Response, error) {
	req, err := http.NewRequest(strings.ToUpper(s.method), s.url, bytes.NewReader([]byte(s.body)))

	if err != nil {
		return nil, fmt.Errorf("request(): %w", err)
	}

	var headers map[string]string

	if err := json.Unmarshal([]byte(s.headers), &headers); err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("request(): %w", err)
		return nil, err
	}

	return resp, nil
}

func (s *Step) prepearRequestData(ctx context.Context) error {
	url, err := s.replaceVars(ctx, s.url, s.findVars(s.url))
	if err != nil {
		return err
	}

	s.url = url

	headers, err := s.replaceVars(ctx, s.headers, s.findVars(s.headers))
	if err != nil {
		return err
	}

	s.headers = headers

	body, err := s.replaceVars(ctx, s.body, s.findVars(s.body))
	if err != nil {
		return err
	}

	s.body = body

	return nil
}

func (s *Step) findVars(text string) []string {
	matches := regexp.
		MustCompile(`\{\{([a-zA-Z0-9_]+)\}\}`).
		FindAllStringSubmatch(text, -1)

	var vars []string
	for _, match := range matches {
		vars = append(vars, match[1])
	}

	return vars
}

func (s *Step) replaceVars(ctx context.Context, text string, vars []string) (string, error) {
	result := text

	for _, name := range vars {
		value := ctx.Value(name)
		if value == nil {
			return result, fmt.Errorf("replaceVars(): variable %q not found in context", name)
		}

		placeholder := "{{" + name + "}}"
		result = strings.ReplaceAll(result, placeholder, fmt.Sprint(value))
	}

	return result, nil
}

func (s *Step) extract(resp *http.Response) (map[string]any, error) {
	vars := make(map[string]any, len(s.extractData))
	for name, path := range s.extractData {
		realPath, ok := path.(string)
		if !ok {
			return nil, fmt.Errorf("extract faild: path is not a string")
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body := string(data)
		result := gjson.Get(body, realPath)

		if !result.Exists() {
			return nil, fmt.Errorf("extract failed: key %s is no exists", path)
		}
		vars[name] = result.Value()
	}

	return vars, nil
}
