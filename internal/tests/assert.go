package tests

import (
	"fmt"
	"reflect"
)

type Assert struct {
	StatusCode int            `json:"statusCode"`
	Equals     map[string]any `json:"equals"`
	Exists     []string       `json:"exists"`
}

func (a *Assert) Check(resp *Response) []error {
	assertErrors := make([]error, 0)

	if a.StatusCode != 0 {
		if err := a.checkStatus(resp); err != nil {
			assertErrors = append(assertErrors, err)
		}
	}

	if a.Equals != nil {
		if errs := a.checkEquals(resp); len(errs) > 0 {
			assertErrors = append(assertErrors, errs...)
		}
	}

	if len(a.Exists) > 0 {
		if errs := a.checkExists(resp); len(errs) > 0 {
			assertErrors = append(assertErrors, errs...)
		}
	}

	return assertErrors
}

func (a *Assert) checkStatus(resp *Response) error {
	if a.StatusCode != resp.StatusCode {
		return fmt.Errorf("Assert status-code - %d !=  Response status-code %d", a.StatusCode, resp.StatusCode)
	}

	return nil
}

func (a *Assert) checkEquals(resp *Response) []error {
	errs := make([]error, 0)
	for path, expected := range a.Equals {
		actual, err := GetByPath(resp.Body, path)
		if err != nil {
			errs = append(errs, fmt.Errorf("equals failed: path %q not found: %w", path, err))
			continue
		}

		if !reflect.DeepEqual(actual, expected) {
			errs = append(errs, fmt.Errorf("equals failed: path %q expected %v, got %v", path, expected, actual))
		}
	}

	return errs
}

func (a *Assert) checkExists(resp *Response) []error {
	errs := make([]error, 0)

	for _, path := range a.Exists {
		_, err := GetByPath(resp.Body, path)
		if err != nil {
			errs = append(errs, fmt.Errorf("exists failed: path %q not found: %w", path, err))
		}
	}

	return errs
}
