package main

import (
	"fmt"
	"integrationstests/internal/report"
	"integrationstests/internal/tests"
	"os"
)

func main() {
	tests, err := tests.LoadFrom("integration.json")
	if err != nil {
		panic(err)
	}

	result := tests.Run()

	if result != nil {
		consoleReporter := &report.ConsoleReporter{
			Writer: os.Stdout,
		}
		consoleReporter.Report(result)

		file, err := os.Create("integration_result.json")
		defer file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		JSONReporter := &report.JSONReporter{
			Writer: file,
		}

		JSONReporter.Report(result)
	}
}
