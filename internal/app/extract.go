package app

import (
	"fmt"

	"github.com/tidwall/gjson"
)

type Extract map[string]string

type ExtractResult struct {
	Success   bool
	Variables map[string]any
	Error     error
}

func (extract *Extract) Do(resp *Response) *ExtractResult {
	success := true
	vars := make(map[string]any)
	var err error

	for name, path := range *extract {
		result := gjson.Get(resp.Body, path)

		if !result.Exists() {
			success = false
			err = fmt.Errorf("failed: extract by key %s is no exists", path)
			break
		}
		vars[name] = result.Value()
	}

	return &ExtractResult{
		Success:   success,
		Variables: vars,
		Error:     err,
	}
}
