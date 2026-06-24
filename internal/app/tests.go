package app

import (
	"encoding/json"
	"os"
)

// IntegrationTests ====>
type IntegrationTests struct {
	BaseURL string     `json:"baseURL"`
	Cases   []TestCase `json:"cases"`
}

type IntegrationTestsResult struct {
	CaseResults []TestCaseResult
}

func (t *IntegrationTests) Run() *IntegrationTestsResult {
	if len(t.Cases) != 0 {
		var result IntegrationTestsResult

		ctx := NewTestContext()

		for _, testCase := range t.Cases {
			caseResult := testCase.Do(ctx, t.BaseURL)
			result.CaseResults = append(result.CaseResults, *caseResult)

			if caseResult.Error != nil {
				break
			}
		}

		return &result
	}

	return nil
}

func LoadFrom(filename string) (IntegrationTests, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return IntegrationTests{}, err
	}

	var tests IntegrationTests
	err = json.Unmarshal(data, &tests)
	return tests, err
}

type TestContext struct {
	variables map[string]any
}

func NewTestContext() *TestContext {
	return &TestContext{
		variables: make(map[string]any),
	}
}

func (c *TestContext) SetManyVars(vars map[string]any) {
	for name, value := range vars {
		c.variables[name] = value
	}
}

func (c *TestContext) GetVarByKey(key string) (any, bool) {
	v, ok := c.variables[key]

	return v, ok
}
