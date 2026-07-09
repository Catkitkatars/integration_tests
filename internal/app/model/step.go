package model

import (
	"context"
	"fmt"
	"integrationstests/internal/app/input"

	"github.com/tidwall/gjson"
)

type StepInterface interface {
	Do(context.Context) (*StepResult, error)
}

type Step struct {
	name        string
	request     *Request
	asserts     *Asserts
	extractData map[string]any
}

type StepResult struct {
	AssertResults []*AssertsResult
}

func NewStep(
	name string,
	request *Request,
	asserts *Asserts,
	extract map[string]any) (*Step, error) {
	return &Step{
		name:        name,
		request:     request,
		asserts:     asserts,
		extractData: extract,
	}, nil
}

func InitStep(data *input.StepData) (*Step, error) {
	if data.Request == nil {
		return nil, fmt.Errorf("unknown step type")
	}

	rq := data.Request
	request := NewRequest(rq.Method, rq.URL, rq.Headers, rq.Body)

	asserts := NewAsserts(data.Asserts.StatusCode, data.Asserts.Equals, data.Asserts.Exists)

	return NewStep(data.Name, request, asserts, data.Extract)
}

func (s *Step) Do(ctx context.Context) (*StepResult, error) {
	err := s.request.FindAndReplaceVars(ctx)
	if err != nil {
		return nil, fmt.Errorf("Step.Do(): %w", err)
	}

	resp, err := s.request.Send()
	if err != nil {
		return nil, fmt.Errorf("Step.Do(): %w", err)
	}

	assertsRes := s.asserts.Check(resp)
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

func (s *Step) extract(resp *Response) (map[string]any, error) {
	vars := make(map[string]any, len(s.extractData))
	for name, path := range s.extractData {
		realPath, ok := path.(string)
		if !ok {
			return nil, fmt.Errorf("extract faild: path is not a string")
		}
		result := gjson.Get(resp.Body, realPath)

		if !result.Exists() {
			return nil, fmt.Errorf("extract failed: key %s is no exists", path)
		}
		vars[name] = result.Value()
	}

	return vars, nil
}
