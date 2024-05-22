package main

import (
	"clams/processor/dynds"
	"database/sql"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	dbConx := getDatabaseConnexion()
	defer dbConx.Close()

	ds := dynds.DynamoDatastore{
		Table:    os.Getenv("WORKSHOP_SIGNUPS_TABLE_NAME"),
		Endpoint: os.Getenv("DYNAMO_ENDPOINT_OVERRIDE"),
		Region:   os.Getenv("AWS_REGION"),
	}
	ds.Init()

	lambdaHandler := handler{
		dbConx:    dbConx,
		datastore: &ds,
	}

	lambda.Start(lambdaHandler.handleRequest)
}

func getDatabaseConnexion() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(fmt.Errorf("opening Postgres Repository connexion at %s: %v", os.Getenv("DATABASE_URL"), err))
	}
	return db
}
