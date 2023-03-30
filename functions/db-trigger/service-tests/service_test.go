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
			ctx.BeforeSuite(steps.setUpDatabaseClient)
			ctx.BeforeSuite(steps.setUpSqsClient)
			ctx.AfterSuite(steps.closeDatabaseConnection)
			ctx.AfterSuite(steps.stopContainers)
		},
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			ctx.Step(`^there are database trigger notifications in the database$`, steps.thereAreDatabaseTriggerNotifications)
			ctx.Step(`^the database trigger is invoked$`, steps.theDbTriggerLambdaIsInvoked)
			ctx.Step(`^the messages are placed on the queue$`, steps.theMessagesArePlacedOnTheQueue)
			ctx.Step(`^the notifications are removed from the database$`, steps.theNotificationsAreRemovedFromTheDatabase)
		},
		Options: &godog.Options{
			StopOnFailure: true,
			Strict:        true,
			Format:        "pretty",
			Paths:         []string{"features"},
			TestingT:      t, // Testing instance that will run subtests.
		},
	}

	if run := suite.Run(); run != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
