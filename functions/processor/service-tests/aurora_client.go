package service_tests

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type AuroraClient struct {
	host   string
	port   int
	dbconx *sql.DB
}

func (a *AuroraClient) connectToDatabase() *sql.DB {
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		a.host, a.port, "hacktivista", "d0ntHackM3", "hacktionlab", "hacktionlab_workshops",
	)

	dbconx, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	return dbconx
}

func (a *AuroraClient) closeDatabaseConnexion() error {
	err := a.dbconx.Close()
	return err
}

func (a *AuroraClient) createDatabaseEntries() {
	statement := `
		insert into "workshops" (id, title)
		values(1, 'My Exciting Workshop on COBOL')
	`
	_, err := a.dbconx.Exec(statement)
	if err != nil {
		panic(err)
	}

	statement = `
		insert into people (id, forename, surname, email)
		values(1, 'Frank', 'Ostrowski', 'frank.o@gfa.de'),(2, 'Grace', 'Hopper', 'g.hopper@codasyl.mil')
	`
	_, err = a.dbconx.Exec(statement)
	if err != nil {
		panic(err)
	}

	statement = `
			insert into roles(id, role_name)
			values(1, 'Facilitator'),(2, 'Attendee')
	`
	_, err = a.dbconx.Exec(statement)
	if err != nil {
		panic(err)
	}

	statement = `
			insert into workshop_signups(people_id, workshop_id, role_id, signed_up_on)
			values(1, 2, 1, now()), (1, 1, 2, now())
	`
	_, err = a.dbconx.Exec(statement)
	if err != nil {
		panic(err)
	}

}
