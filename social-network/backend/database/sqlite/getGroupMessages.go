package sqlite

import (
	"log"
)

type GroupMessage struct {
	GroupId   int    `json:"group_id"`
	UserId    int    `json:"user_id"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Image     []byte `json:"image"`
	CreatedAt string `json:"created_at"`
}

func GetGroupMessages(groupID int) ([]GroupMessage, error) {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer db.Close()

	// get all messages from the group
	rows, err := db.Query("SELECT group_id, user_id, content, image, created_at FROM group_messages WHERE group_id = ?", groupID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	messages := make([]GroupMessage, 0)

	for rows.Next() {
		var message GroupMessage
		err := rows.Scan(&message.GroupId, &message.UserId, &message.Content, &message.Image, &message.CreatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		// get sender name
		var senderName string
		err = db.QueryRow("SELECT fullname FROM users WHERE user_id = ?", message.UserId).Scan(&senderName)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		message.Sender = senderName

		messages = append(messages, message)
	}

	return messages, nil
}
