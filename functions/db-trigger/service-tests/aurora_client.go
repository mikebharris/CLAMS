package service_tests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
)

type DatabaseClient struct {
	host   string
	port   int
	dbConx *sql.DB
}

func (a *DatabaseClient) connectToDatabase() *sql.DB {
	dbConx, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		a.host, a.port, "hacktivista", "d0ntHackM3", "hacktionlab", "hacktionlab_workshops",
	))
	if err != nil {
		panic(err)
	}

	return dbConx
}

func (a *DatabaseClient) closeDatabaseConnexion() error {
	err := a.dbConx.Close()
	return err
}

func (a *DatabaseClient) insertTriggerNotification(n TriggerNotification) {
	msg, _ := json.Marshal(n)
	statement := `insert into trigger_notifications(message) values($1)`
	_, err := a.dbConx.Exec(statement, msg)
	if err != nil {
		panic(err)
	}
}
