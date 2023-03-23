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
	}
	ds.Init()

	lambdaHandler := handler{
		dbConx:    dbConx,
		datastore: &ds,
	}

	lambda.Start(lambdaHandler.handleRequest)
}

func getDatabaseConnexion() *sql.DB {
	connexionString := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=hacktionlab search_path=hacktionlab_workshops sslmode=disable",
		os.Getenv("RDS_HOST"), os.Getenv("RDS_USER"), os.Getenv("RDS_PASSWORD"))
	db, _ := sql.Open("postgres", connexionString)
	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("pinging Postgres Repository connexion at %s: %v", os.Getenv("RDS_HOST"), err))
	}
	return db
}
