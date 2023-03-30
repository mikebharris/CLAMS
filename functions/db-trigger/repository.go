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

func (r *repository) getTriggerNotifications() []notification {
	rows, _ := r.dbConx.Query("select id, message from trigger_notifications")

	var notifications []notification
	for rows.Next() {
		var n notification
		if err := rows.Scan(&n.id, &n.message); err != nil {
			return []notification{}
		} else {
			notifications = append(notifications, n)
		}
	}
	return notifications
}

func (r *repository) deleteTriggerNotification(id int) {
	r.dbConx.Exec(fmt.Sprintf("delete from trigger_notifications where id = %d", id))
}
