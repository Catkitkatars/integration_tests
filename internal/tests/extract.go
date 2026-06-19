package tests

import (
	"fmt"

	"github.com/tidwall/gjson"
)

type Extract map[string]string

type ExtractResult struct {
	Success   bool
	Variables map[string]any
	Failures  []ExtractFailure
}

type ExtractFailure struct {
	Path    string
	Message string
}

func (e *Extract) Do(resp *Response, extract Extract) *ExtractResult {
	success := true
	var vars map[string]any
	var failures []ExtractFailure

	for name, path := range extract {
		result := gjson.Get(resp.Body, path)

		if !result.Exists() {
			fail := ExtractFailure{
				Path:    path,
				Message: fmt.Sprintf("failed: extract by key %w is no exists", path),
			}
			failures = append(failures, fail)
			continue
		}
		vars[name] = result.Value()
	}

	return &ExtractResult{
		Success:   success,
		Variables: vars,
		Failures:  failures,
	}
}
