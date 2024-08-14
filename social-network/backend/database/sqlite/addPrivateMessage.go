package sqlite

import (
	"log"
)

func AddPrivateMessage(senderId, receiverId int, content string, image []byte, createdAt string) error {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	if image != nil {
		_, err = db.Exec("INSERT INTO private_messages (sender_id, receiver_id, content, image, created_at) VALUES (?, ?, ?, ?, ?)", senderId, receiverId, content, image, createdAt)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		_, err = db.Exec("INSERT INTO private_messages (sender_id, receiver_id, content, created_at) VALUES (?, ?, ?, ?)", senderId, receiverId, content, createdAt)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
