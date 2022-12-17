package service_tests

import (
	"github.com/cucumber/godog"
	"testing"
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
			ctx.Step(`^the Attendee Writer is invoked with an attendee record from BAMS to be processed$`, steps.theAttendeeWriterIsInvokedWithANewAttendeeRecord)
			ctx.Step(`^the Attendee Writer is invoked with an updated attendee record from BAMS to be processed$`, steps.theAttendeeWriterIsInvokedWithAnUpdatedAttendeeRecord)
			ctx.Step(`^an attendee record is added to CLAMS$`, steps.theAttendeeIsAddedToTheAttendeesDatastore)
			ctx.Step(`^the attendee record is updated in CLAMS$`, steps.theAttendeeIsUpdatedInTheAttendeesDatastore)
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
