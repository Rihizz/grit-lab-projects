package sqlite

import (
	"log"
)

func AddGroupPost(groupID, userID int, content string, image []byte, createdAt string) error {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	if image != nil {
		_, err = db.Exec("INSERT INTO group_posts (group_id, user_id, content, image, created_at) VALUES (?, ?, ?, ?, ?)", groupID, userID, content, image, createdAt)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		_, err = db.Exec("INSERT INTO group_posts (group_id, user_id, content, created_at) VALUES (?, ?, ?, ?)", groupID, userID, content, createdAt)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
