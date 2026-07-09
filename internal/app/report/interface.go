package report

import "integrationstests/internal/app/model"

type ReporterInterface interface {
	Report(result *model.CaseResult)
}
