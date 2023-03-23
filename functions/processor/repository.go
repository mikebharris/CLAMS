package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type repository struct {
	dbConx *sql.DB
}

func (r *repository) getSignupRecord(signupId int) (WorkshopSignupRecord, error) {
	statement := "select w.title, concat(p.forename, ' ', p.surname) as name, r.role_name " +
		"from people p join workshop_signups s on s.people_id = p.id " +
		"join workshops w on s.workshop_id = w.id " +
		"join roles r on s.role_id = r.id " +
		"where s.id = $1"
	row := r.dbConx.QueryRow(statement, signupId)

	var record WorkshopSignupRecord
	record.WorkshopSignupId = signupId
	if err := row.Scan(&record.WorkshopTitle, &record.Name, &record.Role); err != nil {
		return WorkshopSignupRecord{}, err
	}

	return record, nil
}

func (r *repository) getAttendeeNames(title string) []string {
	statement := "select concat(p.forename, ' ', p.surname) as name " +
		"from people p join workshop_signups s on s.people_id = p.id " +
		"join workshops w on s.workshop_id = w.id " +
		"join roles r on s.role_id = r.id " +
		"where w.title = '$1' and r.role_name = 'attendee'"

	rows, _ := r.dbConx.Query(statement, title)

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return []string{}
		} else {
			names = append(names, name)
		}
	}
	return names
}
