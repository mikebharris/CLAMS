package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

var awsAccountNumber = flag.Uint("account-number", 0, "Account number of AWS deployment target")
var environment = flag.String("environment", "nonprod", "Target environment = prod, nonprod, etc")
var opConfirmed = flag.Bool("confirm", false, "For destructive operations this should be set to true rather than false")
var lambdasToBuild = flag.String("lambda", "all", "Which Lambda functions to test and/or build: <name-of-lambda> or all")
var stage = flag.String("stage", "", "Deployment stage: unit-test, build, int-test, init, plan, apply, destroy")

var lambdas = []string{"attendee-writer", "attendees-api", "db-trigger", "processor"}

func main() {
	flag.Parse()
	switch *stage {
	case "unit-test":
		if *lambdasToBuild == "all" {
			unitTestAll()
		} else {
			unitTest(*lambdasToBuild)
		}
	case "build":
		if *lambdasToBuild == "all" {
			buildAll()
		} else {
			buildLambda(*lambdasToBuild)
		}
	case "int-test":
		if *lambdasToBuild == "all" {
			intTestAll()
		} else {
			intTest(*lambdasToBuild)
		}
	case "init":
		fallthrough
	case "plan":
		fallthrough
	case "apply":
		fallthrough
	case "destroy":
		runTerraformCommandForRegion(*stage, "eu-west-1")
	}
}

func runTerraformCommandForRegion(tfOp string, awsRegion string) {
	tf := setupTerraformExec(context.Background())

	var stdout bytes.Buffer
	tf.SetStdout(&stdout)

	var tfLog strings.Builder
	tf.SetLogger(log.New(&tfLog, "log: ", log.LstdFlags))

	tfWorkingBucket := fmt.Sprintf("%d-%s-terraform-deployments", *awsAccountNumber, awsRegion)
	switch tfOp {
	case "init":
		terraformInit(tf, tfWorkingBucket, awsRegion)
	case "plan":
		terraformPlan(tf, tfWorkingBucket, *awsAccountNumber, *environment, false)
	case "apply":
		if *opConfirmed {
			terraformApply(tf, tfWorkingBucket, *awsAccountNumber, *environment)
		} else {
			log.Println("destructive apply not confirmed running plan instead...")
			terraformPlan(tf, tfWorkingBucket, *awsAccountNumber, *environment, false)
		}
	case "destroy":
		if *opConfirmed {
			terraformDestroy(tf, tfWorkingBucket, *awsAccountNumber, *environment)
		} else {
			log.Println("destructive destroy not confirmed running plan destroy instead...")
			terraformPlan(tf, tfWorkingBucket, *awsAccountNumber, *environment, true)
		}
	case "skip":
	default:
		log.Fatalf("Bad operation: --tfop should be one of init, plan, apply, skip, or destroy")
	}

	fmt.Println("\nterraform log: \n******************\n", tfLog.String())
	fmt.Println("\nterraform stdout: \n******************\n", stdout.String())
}

func setupTerraformExec(ctx context.Context) *tfexec.Terraform {
	log.Println("installing Terraform...")
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.6")),
	}

	execPath, err := installer.Install(ctx)
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	workingDir := "terraform"
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}
	return tf
}

func terraformInit(tf *tfexec.Terraform, tfWorkingBucket string, awsRegion string) {
	log.Println("initialising Terraform...")
	if err := tf.Init(context.Background(),
		tfexec.Upgrade(true),
		tfexec.BackendConfig(fmt.Sprintf("key=tfstate/%s/cmi-schedulers-dsr-automation.json", *environment)),
		tfexec.BackendConfig(fmt.Sprintf("bucket=%s", tfWorkingBucket)),
		tfexec.BackendConfig(fmt.Sprintf("region=%s", awsRegion))); err != nil {
		log.Fatalf("error running Init: %s", err)
	}
}

func terraformPlan(tf *tfexec.Terraform, tfWorkingBucket string, awsAccountNumber uint, environment string, destroyFlag bool) {
	if destroyFlag {
		log.Println("planning Terraform destroy...")
	} else {
		log.Println("planning Terraform apply...")
	}
	_, err := tf.Plan(context.Background(),
		tfexec.Refresh(true),
		tfexec.Destroy(destroyFlag),
		tfexec.Var(fmt.Sprintf("terraform_working_bucket=%s", tfWorkingBucket)),
		tfexec.Var(fmt.Sprintf("account_number=%d", awsAccountNumber)),
		tfexec.Var(fmt.Sprintf("environment=%s", environment)),
		tfexec.VarFile(fmt.Sprintf("environments/%s.tfvars", environment)),
	)
	if err != nil {
		log.Fatalf("error running Plan: %s", err)
	}
}

func terraformApply(tf *tfexec.Terraform, workingBucket string, awsAccountNumber uint, environment string) {
	log.Println("applying Terraform...")
	if err := tf.Apply(context.Background(),
		tfexec.Refresh(true),
		tfexec.Var(fmt.Sprintf("terraform_working_bucket=%s", workingBucket)),
		tfexec.Var(fmt.Sprintf("account_number=%d", awsAccountNumber)),
		tfexec.Var(fmt.Sprintf("environment=%s", environment)),
		tfexec.VarFile(fmt.Sprintf("environments/%s.tfvars", environment)),
	); err != nil {
		log.Fatalf("error running Apply: %s", err)
	}
	displayTerraformOutputs(tf)
}

func terraformDestroy(tf *tfexec.Terraform, workingBucket string, awsAccountNumber uint, environment string) {
	log.Println("destroying all the things...")
	if err := tf.Destroy(context.Background(),
		tfexec.Refresh(true),
		tfexec.Var(fmt.Sprintf("terraform_working_bucket=%s", workingBucket)),
		tfexec.Var(fmt.Sprintf("account_number=%d", awsAccountNumber)),
		tfexec.Var(fmt.Sprintf("environment=%s", environment)),
		tfexec.VarFile(fmt.Sprintf("environments/%s.tfvars", environment)),
	); err != nil {
		log.Fatalf("error running Destroy: %s", err)
	}
	displayTerraformOutputs(tf)
}

func displayTerraformOutputs(tf *tfexec.Terraform) {
	outputs, err := tf.Output(context.Background())
	if err != nil {
		log.Fatalf("Error outputting outputs: %v", err)
	}
	if len(outputs) > 0 {
		fmt.Println("Terraform outputs:")
	}
	for key := range outputs {
		if outputs[key].Sensitive {
			continue
		}
		fmt.Println(fmt.Sprintf("%s = %s\n", key, string(outputs[key].Value)))
	}
}

func buildLambda(lambdaName string) {
	unitTest(lambdaName)
	buildTarget(lambdaName)
	intTest(lambdaName)
}

func unitTestAll() {
	for _, lambda := range lambdas {
		unitTest(lambda)
	}
}

func buildAll() {
	for _, lambda := range lambdas {
		buildTarget(lambda)
	}
}
func intTestAll() {
	for _, lambda := range lambdas {
		intTest(lambda)
	}
}

func unitTest(lambdaName string) {
	log.Printf("running tests for %s Lambda...\n", lambdaName)
	stdout := runCmdIn(fmt.Sprintf("lambdas/%s", lambdaName), "make", "unit-test")
	fmt.Println("unit tests passed; stdout = ", stdout)
}

func buildTarget(lambdaName string) {
	log.Printf("building %s Lambda...\n", lambdaName)
	stdout := runCmdIn(fmt.Sprintf("lambdas/%s", lambdaName), "make", "target")
	fmt.Println("build succeeded; stdout = ", stdout)
}

func intTest(lambdaName string) {
	log.Printf("running integration tests for %s Lambda...\n", lambdaName)
	stdout := runCmdIn(fmt.Sprintf("lambdas/%s", lambdaName), "make", "int-test")
	fmt.Println("integration tests passed; stdout = ", stdout)
}

func runCmdIn(dir string, command string, args ...string) string {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	var stdout strings.Builder
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		log.Fatalf("error running %s %s: %s\n\n******************\n\n%s\n******************\n\n", command, args, err, stdout.String())
	}
	return stdout.String()
}
