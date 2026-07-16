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

// Нет смысла парсить body в map[string]any
type RequestData struct {
	Method  string `json:"method"`
	URL     string `json:"url"`
	Headers string `json:"headers"`
	Body    string `json:"body"`
}

// Попробовать AssertsData map[string]any и парсить уже внутри функции конструктора

type AssertsData struct {
	StatusCode int            `json:"statusCode"`
	Equals     map[string]any `json:"equals"`
	Exists     []string       `json:"exists"`
}
