package model

import (
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/tidwall/gjson"
)

const (
	StatusCodeCheckType = "statusCode"
	EqualsCheckType     = "equals"
	ExistsCheckType     = "exists"
)

type Asserts struct {
	statusCode int
	equals     map[string]any
	exists     []string
}

func NewAsserts(
	statusCode int,
	equals map[string]any,
	exists []string,
) *Asserts {
	return &Asserts{
		statusCode: statusCode,
		equals:     equals,
		exists:     exists,
	}
}

type AssertsResult struct {
	Type     string `json:"type"`
	Path     string `json:"path"`
	Expected any    `json:"expected"`
	Actual   any    `json:"actual"`
	Message  string `json:"message"`
}

func (a *Asserts) Check(resp *http.Response) ([]*AssertsResult, error) {
	var failures []*AssertsResult

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Asserts.Check(): error read resp.Body")
	}
	defer resp.Body.Close()

	body := string(data)

	if a.statusCode != 0 {
		if statusFailures := a.checkStatus(resp.StatusCode); len(statusFailures) > 0 {
			failures = append(failures, statusFailures...)
		}
	}

	if a.equals != nil {
		if equalsFailures := a.checkEquals(body); len(equalsFailures) > 0 {
			failures = append(failures, equalsFailures...)
		}
	}

	if len(a.exists) > 0 {
		if existsFailures := a.checkExists(body); len(existsFailures) > 0 {
			failures = append(failures, existsFailures...)
		}
	}

	return failures, nil
}

func (a *Asserts) checkStatus(statusCode int) []*AssertsResult {
	var failures []*AssertsResult
	if a.statusCode != statusCode {
		fail := AssertsResult{
			Type:     StatusCodeCheckType,
			Path:     "statusCode",
			Expected: a.statusCode,
			Actual:   statusCode,
			Message: fmt.Sprintf(
				"asserts failed: status-code - %d != Response status-code %d",
				a.statusCode,
				statusCode,
			),
		}
		failures = append(failures, &fail)
	}

	return failures
}

func (a *Asserts) checkEquals(body string) []*AssertsResult {
	var failures []*AssertsResult

	for path, expected := range a.equals {
		result := gjson.Get(body, path)

		if !result.Exists() {
			fail := AssertsResult{
				Type:     EqualsCheckType,
				Path:     path,
				Expected: expected,
				Actual:   result.Value(),
				Message:  fmt.Sprintf("asserts failed: path %q not found in body", path),
			}
			failures = append(failures, &fail)
			continue
		}

		if !reflect.DeepEqual(result.Value(), expected) {
			fail := AssertsResult{
				Type:     EqualsCheckType,
				Path:     path,
				Expected: expected,
				Actual:   result.Value(),
				Message:  fmt.Sprintf("asserts failed: path %q expected %v, got %v", path, expected, result.Value()),
			}

			failures = append(failures, &fail)
		}
	}

	return failures
}

func (a *Asserts) checkExists(body string) []*AssertsResult {
	var failures []*AssertsResult

	for _, path := range a.exists {
		if value := gjson.Get(body, path); !value.Exists() {
			fail := AssertsResult{
				Type:     ExistsCheckType,
				Path:     path,
				Expected: nil,
				Actual:   nil,
				Message:  fmt.Sprintf("asserts failed: path %q not found", path),
			}

			failures = append(failures, &fail)
		}
	}

	return failures
}
