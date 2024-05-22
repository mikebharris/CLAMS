package integration_tests

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
)

const facilitatorRoleId = 1
const attendeeRoleId = 2

type AuroraClient struct {
	host   string
	port   int
	dbConx *sql.DB
}

func (a *AuroraClient) connectToDatabase() *sql.DB {
	dbConx, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		a.host, a.port, "hacktivista", "d0ntHackM3", "hacktionlab", "hacktionlab_workshops",
	))
	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}

	return dbConx
}

func (a *AuroraClient) closeDatabaseConnexion() error {
	err := a.dbConx.Close()
	return err
}

func (a *AuroraClient) createWorkshopSignup(workshopId, personId, roleId int) int {
	id := rand.Intn(100)
	statement := `
			insert into workshop_signups(id, people_id, workshop_id, role_id, signed_up_on)
			values($1, $2, $3, $4, now())
	`
	_, err := a.dbConx.Exec(statement, id, personId, workshopId, roleId)
	if err != nil {
		log.Fatalf("creating workshop signup: %v", err)
	}
	return id
}

func (a *AuroraClient) createRoles() {
	statement := `
			insert into roles(id, role_name)
			values($1, 'facilitator'),($2, 'attendee')
	`
	_, err := a.dbConx.Exec(statement, facilitatorRoleId, attendeeRoleId)
	if err != nil {
		log.Fatalf("creating roles: %v", err)
	}
}

func (a *AuroraClient) createPerson(forename, surname, email string) int {
	id := rand.Intn(100)

	statement := `
		insert into people(id, forename, surname, email)
		values($1, $2, $3, $4)
	`
	_, err := a.dbConx.Exec(statement, id, forename, surname, email)
	if err != nil {
		log.Fatalf("creating person: %v", err)
	}
	return id
}

func (a *AuroraClient) createWorkshop(title string) int {
	id := rand.Intn(100)
	statement := `
		insert into workshops(id, title)
		values($1, $2)
	`
	_, err := a.dbConx.Exec(statement, id, title)
	if err != nil {
		log.Fatalf("creating workshop: %v", err)
	}
	return id
}
