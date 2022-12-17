package service_tests

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	var steps steps
	steps.t = t

	suite := godog.TestSuite{
		TestSuiteInitializer: func(ctx *godog.TestSuiteContext) {
			ctx.BeforeSuite(steps.startContainers)
			ctx.BeforeSuite(steps.setUpDynamoClient)
			ctx.AfterSuite(steps.stopContainers)
		},
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			ctx.Step(`^some attendee records exist in CLAMS$`, steps.someAttendeeRecordsExistInTheAttendeesDatastore)
			ctx.Step(`^the front-end requests a specific attendee record from the endpoint$`, steps.theFrontendRequestsASpecificRecordFromTheEndpoint)
			ctx.Step(`^the record is returned$`, steps.aSingleRecordIsReturned)

			ctx.Step(`^the front-end requests all records from the endpoint$`, steps.theFrontendRequestsAllRecordsFromTheEndpoint)
			ctx.Step(`^all available records are returned$`, steps.theRecordsAreReturned)

			ctx.Step(`^the front-end requests the stats from the report endpoint$`, steps.theFrontEndRequestsTheStatsFromTheReportEndpoint)
			ctx.Step(`^some statistics about the event are returned$`, steps.someStatisticsAboutTheEventAreReturned)

		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
