package main

import "integrationstests/internal/tests"

func main() {
	Run("integration.json")
}

func Run(filepath string) {
	tests, err := tests.LoadTests(filepath)
	if err != nil {
		panic(err)
	}

	tests.Start()
}
