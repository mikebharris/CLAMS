package main

import (
	"clams/awscfg"
	"database/sql"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	dbConx := getDatabaseConnexion()
	defer dbConx.Close()

	awsConfig := awscfg.GetAwsConfig(sqs.ServiceID, os.Getenv("SQS_ENDPOINT_OVERRIDE"), os.Getenv("AWS_REGION"))

	lambdaHandler := handler{
		dbConx:     dbConx,
		sqsService: sqs.NewFromConfig(*awsConfig),
	}

	lambda.Start(lambdaHandler.handleRequest)
}

func getDatabaseConnexion() *sql.DB {
	dbConx, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(fmt.Errorf("opening Postgres Repository connexion at %s: %v", os.Getenv("DATABASE_URL"), err))
	}
	return dbConx
}
