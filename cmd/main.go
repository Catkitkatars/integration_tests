package main

import (
	"fmt"
	"integrationstests/internal/app"
	"integrationstests/internal/report"
	"os"
)

func main() {
	tests, err := app.LoadFrom("integration.json")
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
