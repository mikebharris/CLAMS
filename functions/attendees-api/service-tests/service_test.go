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

			ctx.Step(`^an attendee record exists in the attendees datastore$`, steps.anAttendeeRecordExistsInTheAttendeesDatastore)
			ctx.Step(`^the front-end requests the attendee record from the API$`, steps.theFrontendFetchesTheRecordFromTheAPI)
			ctx.Step(`^the record is returned$`, steps.theRecordIsReturned)
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
