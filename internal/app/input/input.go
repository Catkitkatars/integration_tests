package input

import (
	"encoding/json"
	"os"
)

type CaseData struct {
	Desc  string         `json:"desc"`
	Vars  map[string]any `json:"vars"`
	Steps []*StepData    `json:"steps"`
}

func LoadFrom(filename string) (*CaseData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var caseData CaseData
	err = json.Unmarshal(data, &caseData)
	return &caseData, err
}

type StepData struct {
	Name    string         `json:"name"`
	Request *RequestData   `json:"rq"`
	Extract map[string]any `json:"extract"`
	Asserts *AssertsData   `json:"asserts"`
}

type RequestData struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    map[string]any    `json:"body"`
}

type AssertsData struct {
	StatusCode int            `json:"statusCode"`
	Equals     map[string]any `json:"equals"`
	Exists     []string       `json:"exists"`
}
