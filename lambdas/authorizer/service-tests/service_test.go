package main

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
			ctx.AfterSuite(steps.stopContainers)
		},
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			ctx.Step(`^The lambda is invoked with valid credentials$`, steps.theLambdaIsInvokedWithValidCredentials)
			ctx.Step(`^A Success response is returned$`, steps.aSuccessResponseIsReturned)
			ctx.Step(`^The lambda is invoked with invalid credentials$`, steps.theLambdaIsInvokedWithInvalidCredentials)
			ctx.Step(`^A Failure response is returned$`, steps.aFailureResponseIsReturned)
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
