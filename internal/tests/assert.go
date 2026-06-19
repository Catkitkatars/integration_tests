package tests

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

type Assert struct {
	StatusCode int            `json:"statusCode"`
	Equals     map[string]any `json:"equals"`
	Exists     []string       `json:"exists"`
}

type AssertResult struct {
	Success  bool            `json:"success"`
	Failures []AssertFailure `json:"failures"`
}

type AssertFailure struct {
	Type     string `json:"type"`
	Path     string `json:"path"`
	Expected any    `json:"expected"`
	Actual   any    `json:"actual"`
	Message  string `json:"message"`
}

func (a *Assert) Check(resp *Response) *AssertResult {
	success := true
	var result []AssertFailure

	if a.StatusCode != 0 {
		if failures := a.checkStatus(resp); failures != nil {
			result = append(result, failures...)
			success = false
		}
	}

	if a.Equals != nil {
		if failures := a.checkEquals(resp); failures != nil {
			result = append(result, failures...)
			success = false
		}
	}

	if len(a.Exists) > 0 {
		if failures := a.checkExists(resp); failures != nil {
			result = append(result, failures...)
			success = false
		}
	}

	return &AssertResult{
		Success:  success,
		Failures: result,
	}
}

func (a *Assert) checkStatus(resp *Response) []AssertFailure {
	var failures []AssertFailure
	if a.StatusCode != resp.StatusCode {
		fail := AssertFailure{
			Type:     StatusCodeCheckType,
			Path:     "statusCode",
			Expected: a.StatusCode,
			Actual:   resp.StatusCode,
			Message: fmt.Sprintf(
				"failed: status-code - %d != Response status-code %d",
				a.StatusCode,
				resp.StatusCode,
			),
		}
		failures = append(failures, fail)
	}

	return failures
}

func (a *Assert) checkEquals(resp *Response) []AssertFailure {
	var failures []AssertFailure

	for path, expected := range a.Equals {
		result := gjson.Get(resp.Body, path)

		if !result.Exists() {
			fail := AssertFailure{
				Type:     EqualsCheckType,
				Path:     path,
				Expected: expected,
				Actual:   result.Value(),
				Message:  fmt.Sprintf("failed: path %q not found in body", path),
			}
			failures = append(failures, fail)
			continue
		}

		if !reflect.DeepEqual(result.Value(), expected) {
			fail := AssertFailure{
				Type:     EqualsCheckType,
				Path:     path,
				Expected: expected,
				Actual:   result.Value(),
				Message:  fmt.Sprintf("failed: path %q expected %v, got %v", path, expected, result.Value()),
			}

			failures = append(failures, fail)
		}
	}

	return failures
}

func (a *Assert) checkExists(resp *Response) []AssertFailure {
	var failures []AssertFailure

	for _, path := range a.Exists {
		if value := gjson.Get(resp.Body, path); !value.Exists() {
			fail := AssertFailure{
				Type:     ExistsCheckType,
				Path:     path,
				Expected: nil,
				Actual:   nil,
				Message:  fmt.Sprintf("failed: path %q not found", path),
			}

			failures = append(failures, fail)
		}
	}

	return failures
}
