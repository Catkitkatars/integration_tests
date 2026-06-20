package report

import "integrationstests/internal/tests"

type ReporterInterface interface {
	Report(result *tests.IntegrationTestsResult)
}
