package sqlite

import (
	"database/sql"
	"log"
)

type Notification struct {
	NotifId  int
	UserId    int
	SenderId  int
	Content   string
	Read bool
	Type string
	Date      string
	Groupid   int
}


func GetNotifications(UserID int) ([]Notification, error) {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT id, user_id, sender_id, content, is_read, type, created_at, reference_id FROM notifications WHERE user_id = ?", UserID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	var notifications []Notification

	for rows.Next() {
		var notification Notification
		var nullableGroupId sql.NullInt64
		err := rows.Scan(&notification.NotifId, &notification.UserId, &notification.SenderId, &notification.Content, &notification.Read, &notification.Type, &notification.Date, &nullableGroupId)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		// Check if nullableGroupId is valid, if so convert to int and assign to Groupid
		if nullableGroupId.Valid {
			notification.Groupid = int(nullableGroupId.Int64)
		} else {
			notification.Groupid = 0 // or any other default value
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}