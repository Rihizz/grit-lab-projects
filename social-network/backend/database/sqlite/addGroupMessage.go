package sqlite

import (
	"log"
)

func AddGroupMessage(groupId, senderId int, content string, image []byte, time string) {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return
	}

	defer db.Close()

	if image != nil {
		_, err = db.Exec("INSERT INTO group_messages (group_id, user_id, content, image, created_at) VALUES (?, ?, ?, ?, ?)", groupId, senderId, content, image, time)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		_, err = db.Exec("INSERT INTO group_messages (group_id, user_id, content, created_at) VALUES (?, ?, ?, ?)", groupId, senderId, content, time)
		if err != nil {
			log.Println(err)
			return
		}
	}

	// // add message to the database
	// _, err = db.Exec("INSERT INTO group_messages (group_id, user_id, content, created_at) VALUES (?, ?, ?, ?)", groupId, senderId, content, time)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
}
