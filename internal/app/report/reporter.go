package report

import (
	"encoding/json"
	"fmt"
	"integrationstests/internal/app/model"
	"io"
)

type ConsoleReporter struct {
	Writer io.Writer
}

func (r *ConsoleReporter) Report(result *model.CaseResult) {
	for _, stepResult := range result.StepResults {

		if len(stepResult.AssertResults) > 0 {
			fmt.Fprintln(r.Writer, "  FAIL")
			for _, failure := range stepResult.AssertResults {
				fmt.Fprintf(r.Writer, "    - %s\n", failure.Message)
			}
			continue
		}

		fmt.Fprintln(r.Writer, "  PASS")
	}
}

type JSONReporter struct {
	Writer io.Writer
}

func (r *JSONReporter) Report(result *model.CaseResult) {
	encoder := json.NewEncoder(r.Writer)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(result); err != nil {
		fmt.Printf("encode json report: %v\n", err)
	}
}
