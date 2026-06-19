package tests

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

func (t *IntegrationTests) Run() {
	if len(t.Cases) != 0 {
		ctx := &TestContext{
			Variables: make(map[string]any),
		}

		for _, testCase := range t.Cases {
			result := testCase.Do(ctx, t.BaseURL)
		}
	}
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
	Variables map[string]any
}

func NewTestContext() *TestContext {
	return &TestContext{
		Variables: make(map[string]any),
	}
}

func (c *TestContext) SetMany(vars map[string]any) {
	for name, value := range vars {
		c.Variables[name] = value
	}
}
