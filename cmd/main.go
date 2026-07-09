package main

import (
	"fmt"
	"integrationstests/internal/app/input"
	"integrationstests/internal/app/model"
	"integrationstests/internal/app/report"
	"os"
)

func main() {
	data, err := input.LoadFrom("integration.json")
	if err != nil {
		panic(err)
	}

	testCase, err := model.InitCase(data)
	if err != nil {
		panic(err)
	}

	result := testCase.Do()

	if result != nil {
		consoleReporter := &report.ConsoleReporter{
			Writer: os.Stdout,
		}
		consoleReporter.Report(result)

		file, err := os.Create("integration_result.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		JSONReporter := &report.JSONReporter{
			Writer: file,
		}

		JSONReporter.Report(result)
	}
}
