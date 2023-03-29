package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type repository struct {
	dbConx *sql.DB
}

func (r *repository) getTriggerNotifications() []string {
	statement := "select message from trigger_notifications"
	rows, _ := r.dbConx.Query(statement)

	var notifications []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return []string{}
		} else {
			notifications = append(notifications, n)
		}
	}
	return notifications
}
