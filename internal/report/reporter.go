package report

import (
	"encoding/json"
	"fmt"
	"integrationstests/internal/app"
	"io"
)

type ConsoleReporter struct {
	Writer io.Writer
}

func (r *ConsoleReporter) Report(result *app.IntegrationTestsResult) {
	for _, caseResult := range result.CaseResults {
		fmt.Fprintf(r.Writer, "Case: %s\n", caseResult.Name)
		if caseResult.Error != nil {
			fmt.Fprintf(r.Writer, "  ERROR: %v\n", caseResult.Error)
			continue
		}

		if caseResult.AssertResult != nil && !caseResult.AssertResult.Success {
			fmt.Fprintln(r.Writer, "  FAIL")
			for _, failure := range caseResult.AssertResult.Failures {
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

func (r *JSONReporter) Report(result *app.IntegrationTestsResult) {
	encoder := json.NewEncoder(r.Writer)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(result); err != nil {
		fmt.Printf("encode json report: %v\n", err)
	}
}
