package model

import (
	"context"
	"fmt"
	"integrationstests/internal/app/input"
)

type Case struct {
	vars  map[string]any
	steps []*Step
}

type CaseResult struct {
	StepResults []*StepResult
	Error       error
}

func NewCase(vars map[string]any, steps []*Step) *Case {
	return &Case{
		vars:  vars,
		steps: steps,
	}
}

func InitCase(data *input.CaseData) (*Case, error) {
	vars := data.Vars
	if vars == nil {
		vars = make(map[string]any)
	}

	steps := make([]*Step, 0, len(data.Steps))

	for _, s := range data.Steps {
		step, err := InitStep(s)
		if err != nil {
			return nil, fmt.Errorf("InitStep: %w", err)
		}

		steps = append(steps, step)
	}

	return NewCase(vars, steps), nil
}

func (c *Case) Do() *CaseResult {
	stepResults := make([]*StepResult, len(c.steps))
	var stepError error
	ctx := context.Background()
	for _, step := range c.steps {
		stepRes, err := step.Do(ctx)

		if err != nil {
			stepError = err
			break
		}

		stepResults = append(stepResults, stepRes)
	}

	return &CaseResult{
		stepResults,
		stepError,
	}
}
