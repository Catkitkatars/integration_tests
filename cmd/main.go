package main

import "integrationstests/internal/tests"

func main() {
	tests, err := tests.LoadFrom("integration.json")
	if err != nil {
		panic(err)
	}

	tests.Run()
}
