package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type repository struct {
	dbConx *sql.DB
}

type notification struct {
	id      int
	message string
}

func (r *repository) getTriggerNotifications() ([]notification, error) {
	rows, err := r.dbConx.Query("select id, message from trigger_notifications")
	if err != nil {
		return []notification{}, fmt.Errorf("fetching trigger notifications: %v", err)
	}

	var notifications []notification
	for rows.Next() {
		var n notification
		if err := rows.Scan(&n.id, &n.message); err != nil {
			return []notification{}, nil
		} else {
			notifications = append(notifications, n)
		}
	}
	return notifications, nil
}

func (r *repository) deleteTriggerNotification(id int) error {
	_, err := r.dbConx.Exec(fmt.Sprintf("delete from trigger_notifications where id = %d", id))
	return err
}
