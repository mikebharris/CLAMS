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
			ctx.BeforeSuite(steps.setUpAuroraClient)
			ctx.AfterSuite(steps.stopContainers)
		},
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			ctx.Step(`^workshop signup records exist in the database$`, steps.theWorkshopSignupRequestExistsInTheDatabase)
			ctx.Step(`^the workshop signup processor receives a notification$`, steps.theProcessorLambdaIsInvoked)
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
