package app

import (
	"fmt"
	"regexp"
	"strings"
)

var varPattern = regexp.MustCompile(`\{\{([a-zA-Z0-9_]+)\}\}`)

func FindAndReplaceVars(ctx *TestContext, req *Request) error {
	url, err := replaceVars(ctx, req.URL, findVars(req.URL))
	if err != nil {
		return err
	}
	req.URL = url

	headers := make(map[string]string, len(req.Headers))
	for key, value := range req.Headers {
		header, err := replaceVars(ctx, value, findVars(value))
		if err != nil {
			return err
		}

		headers[key] = header
	}
	req.Headers = headers

	body, err := replaceVarsFromBody(ctx, req.Body)
	if err != nil {
		return err
	}

	req.Body = body

	return nil
}

func replaceVarsFromBody(ctx *TestContext, body map[string]any) (map[string]any, error) {
	handledBody := make(map[string]any, len(body))

	for key, value := range body {
		handledValue, err := replaceVarsFromValue(ctx, value)
		if err != nil {
			return nil, err
		}

		handledBody[key] = handledValue
	}

	return handledBody, nil
}

func replaceVarsFromValue(ctx *TestContext, value any) (any, error) {
	switch value := value.(type) {
	case string:
		return replaceVars(ctx, value, findVars(value))

	case map[string]any:
		return replaceVarsFromBody(ctx, value)

	case []any:
		handledSlice := make([]any, 0, len(value))

		for _, item := range value {
			handledItem, err := replaceVarsFromValue(ctx, item)
			if err != nil {
				return nil, err
			}

			handledSlice = append(handledSlice, handledItem)
		}

		return handledSlice, nil

	default:
		return value, nil
	}
}

func findVars(text string) []string {
	matches := varPattern.FindAllStringSubmatch(text, -1)

	var vars []string
	for _, match := range matches {
		vars = append(vars, match[1])
	}

	return vars
}

func replaceVars(ctx *TestContext, text string, vars []string) (string, error) {
	result := text

	for _, name := range vars {
		value, ok := ctx.GetVarByKey(name)
		if !ok {
			return result, fmt.Errorf("variable %q not found in test context", name)
		}

		placeholder := "{{" + name + "}}"
		result = strings.ReplaceAll(result, placeholder, fmt.Sprint(value))
	}

	return result, nil
}

func interpolateVars(ctx *TestContext, text string) (string, error) {
	return replaceVars(ctx, text, findVars(text))
}
