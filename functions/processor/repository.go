package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type repository struct {
	dbConx *sql.DB
}

func (r *repository) getSignupRecord(signupId int) (WorkshopSignupRecord, error) {
	statement := `select workshop_id from workshop_signups where id = $1`
	row := r.dbConx.QueryRow(statement, signupId)

	var workshopId int
	if err := row.Scan(&workshopId); err != nil {
		return WorkshopSignupRecord{}, err
	}

	statement = `select id, title from workshops where id = $1`
	row = r.dbConx.QueryRow(statement, signupId)

	var record WorkshopSignupRecord
	if err := row.Scan(&record.WorkshopId, &record.WorkshopTitle); err != nil {
		return WorkshopSignupRecord{}, err
	}

	return record, nil
}
