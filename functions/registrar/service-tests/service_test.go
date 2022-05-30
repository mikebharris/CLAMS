package service_tests

import (
	"github.com/cucumber/godog"
	"testing"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			var steps steps
			steps.t = t

			ctx.Before(steps.startContainers)
			ctx.Before(steps.setUpDynamoClient)
			ctx.After(steps.stopContainers)

			ctx.Step(`^the Registrar is invoked with an attendee record from BAMS to be processed$`, steps.theRegistrarIsInvokedWithANewAttendeeRecord)
			ctx.Step(`^the Registrar is invoked with an updated attendee record from BAMS to be processed$`, steps.theRegistrarIsInvokedWithAnUpdatedAttendeeRecord)

			ctx.Step(`^the attendee is added to the Attendees Datastore$`, steps.theAttendeeIsAddedToTheAttendeesDatastore)
			ctx.Step(`^the attendee is updated in the Attendees Datastore$`, steps.theAttendeeIsUpdatedInTheAttendeesDatastore)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if run := suite.Run(); run != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
