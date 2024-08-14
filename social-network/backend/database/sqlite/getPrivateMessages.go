package sqlite

import (
	"encoding/base64"
	"log"
)

type PrivateMessage struct {
	SenderID     int    `json:"sender_id"`
	ReceiverID   int    `json:"receiver_id"`
	SenderName   string `json:"sender_name"`
	ReceiverName string `json:"receiver_name"`
	Content      string `json:"content"`
	Image        string `json:"image"`
	CreatedAt    string `json:"created_at"`
}

func GetPrivateMessages(senderID, receiverID int) ([]PrivateMessage, error) {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer db.Close()
	log.Println("senderID: ", senderID)
	log.Println("receiverID: ", receiverID)
	// get all messages from the group
	rows, err := db.Query("SELECT sender_id, receiver_id, content, image, created_at FROM private_messages WHERE (receiver_id = ? AND sender_id = ?) OR (sender_id = ? AND receiver_id = ?) ORDER BY id DESC", receiverID, senderID, receiverID, senderID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	messages := make([]PrivateMessage, 0)

	for rows.Next() {
		var message PrivateMessage
		var imageData []byte
		err := rows.Scan(&message.SenderID, &message.ReceiverID, &message.Content, &imageData, &message.CreatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		// get sender name and receiver name
		var senderName string
		var receiverName string
		err = db.QueryRow("SELECT fullname FROM users WHERE user_id = ?", message.SenderID).Scan(&senderName)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		err = db.QueryRow("SELECT fullname FROM users WHERE user_id = ?", message.ReceiverID).Scan(&receiverName)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		// convert image to base64 string
		if imageData != nil {
			message.Image = base64.StdEncoding.EncodeToString(imageData)
		}

		message.SenderName = senderName
		message.ReceiverName = receiverName

		messages = append(messages, message)
	}

	return messages, nil

}
