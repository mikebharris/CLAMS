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

	awsConfig := awscfg.GetAwsConfig(sqs.ServiceID, os.Getenv("SQS_ENDPOINT_OVERRIDE"))

	lambdaHandler := handler{
		dbConx:     dbConx,
		sqsService: sqs.NewFromConfig(*awsConfig),
	}

	lambda.Start(lambdaHandler.handleRequest)
}

func getDatabaseConnexion() *sql.DB {
	connexionString := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s search_path=hacktionlab_workshops sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, _ := sql.Open("postgres", connexionString)
	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("pinging Postgres Repository connexion at %s: %v", os.Getenv("DB_HOST"), err))
	}
	return db
}
