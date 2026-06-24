package report

import "integrationstests/internal/app"

type ReporterInterface interface {
	Report(result *app.IntegrationTestsResult)
}
