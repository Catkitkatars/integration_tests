package model

import (
	"fmt"
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

type AssertFailure struct {
}

func (a *Asserts) Check(resp *Response) []*AssertsResult {
	var failures []*AssertsResult

	if a.statusCode != 0 {
		if statusFailures := a.checkStatus(resp); len(statusFailures) > 0 {
			failures = append(failures, statusFailures...)
		}
	}

	if a.equals != nil {
		if equalsFailures := a.checkEquals(resp); len(equalsFailures) > 0 {
			failures = append(failures, equalsFailures...)
		}
	}

	if len(a.exists) > 0 {
		if existsFailures := a.checkExists(resp); len(existsFailures) > 0 {
			failures = append(failures, existsFailures...)
		}
	}

	return failures
}

func (a *Asserts) checkStatus(resp *Response) []*AssertsResult {
	var failures []*AssertsResult
	if a.statusCode != resp.StatusCode {
		fail := AssertsResult{
			Type:     StatusCodeCheckType,
			Path:     "statusCode",
			Expected: a.statusCode,
			Actual:   resp.StatusCode,
			Message: fmt.Sprintf(
				"asserts failed: status-code - %d != Response status-code %d",
				a.statusCode,
				resp.StatusCode,
			),
		}
		failures = append(failures, &fail)
	}

	return failures
}

func (a *Asserts) checkEquals(resp *Response) []*AssertsResult {
	var failures []*AssertsResult

	for path, expected := range a.equals {
		result := gjson.Get(resp.Body, path)

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

func (a *Asserts) checkExists(resp *Response) []*AssertsResult {
	var failures []*AssertsResult

	for _, path := range a.exists {
		if value := gjson.Get(resp.Body, path); !value.Exists() {
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
