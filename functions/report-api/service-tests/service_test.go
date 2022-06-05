package service_tests

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			var steps steps
			steps.t = t

			ctx.Before(steps.startContainers)
			ctx.Before(steps.setUpDynamoClient)
			ctx.After(steps.stopContainers)

			ctx.Step(`^some attendee records exist in the attendees datastore$`, steps.someAttendeeRecordsExistInTheAttendeesDatastore)
			ctx.Step(`^the front-end requests the stats from the report API$`, steps.theFrontEndRequestsTheStatsFromTheReportAPI)
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
