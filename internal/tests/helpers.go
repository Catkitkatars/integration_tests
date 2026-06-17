package tests

import (
	"fmt"
	"os"
	"strings"
)

func GetByPath(body map[string]any, path string) (any, error) {
	parts := strings.Split(path, ".")

	var current any = body

	for _, part := range parts {
		currentMap, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("failed: path %q is not an object at %q", path, part)
		}

		value, ok := currentMap[part]
		if !ok {
			return nil, fmt.Errorf("failed: path %q not found at %q", path, part)
		}

		current = value
	}
	return current, nil
}

func PrintError(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "\033[31mERROR:\033[0m %v\n", err)
}

func PrintDone() {
	fmt.Fprintln(os.Stdout, "\033[32mPASS: all good\033[0m")
}
