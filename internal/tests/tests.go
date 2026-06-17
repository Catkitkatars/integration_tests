package tests

import (
	"encoding/json"
	"fmt"
	"os"
)

// IntegrationTests ====>
type IntegrationTests struct {
	BaseURL string     `json:"baseURL"`
	Cases   []TestCase `json:"cases"`
}

func (t *IntegrationTests) Start() {
	if len(t.Cases) != 0 {
		ctx := &TestContext{
			Variables: make(Variables),
		}

		for _, testCase := range t.Cases {
			fmt.Printf("Started case - %s\n", testCase.Name)
			errs := testCase.Exec(ctx, t.BaseURL)

			if len(errs) > 0 {
				for _, err := range errs {
					PrintError(err)
				}
				continue
			}

			PrintDone()
		}
	}
}

func LoadTests(filename string) (IntegrationTests, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return IntegrationTests{}, err
	}

	var tests IntegrationTests
	err = json.Unmarshal(data, &tests)
	return tests, err
}

// TestCase ====>
type TestCase struct {
	Name    string  `json:"name"`
	Request Request `json:"request"`
	Extract Extract `json:"extract"`
	Asserts Assert  `json:"asserts"`
}

func (t *TestCase) Exec(ctx *TestContext, baseURL string) []error {
	url := baseURL + t.Request.URL
	httpResp, respTime, err := t.Request.Send(ctx, url)

	testErrors := make([]error, 0)

	if err != nil {
		testErrors = append(testErrors, err)
		return testErrors
	}

	resp, err := NewResponse(httpResp, respTime)

	if err != nil {
		testErrors = append(testErrors, err)
		return testErrors
	}

	vars, err := NewVariables(resp, t.Extract)

	if err != nil {
		testErrors = append(testErrors, err)
		return testErrors
	}

	ctx.SetMany(vars)
	return t.Asserts.Check(resp)
}

// TestContext ====>
type TestContext struct {
	Variables Variables
}

func NewTestContext() *TestContext {
	return &TestContext{
		Variables: make(Variables),
	}
}

func (c *TestContext) Set(name string, value any) {
	c.Variables[name] = value
}

func (c *TestContext) SetMany(vars Variables) {
	for name, value := range vars {
		c.Variables[name] = value
	}
}

func (c *TestContext) Get(name string) (any, bool) {
	value, ok := c.Variables[name]
	return value, ok
}

// Extract ====>
type Extract map[string]string

// Variables ====>
type Variables map[string]any

func NewVariables(resp *Response, extract Extract) (Variables, error) {
	vars := make(Variables)

	for name, path := range extract {
		current, err := GetByPath(resp.Body, path)

		if err != nil {
			return nil, fmt.Errorf("extract %q failed: %w", name, err)
		}
		vars[name] = current
	}

	return vars, nil
}
