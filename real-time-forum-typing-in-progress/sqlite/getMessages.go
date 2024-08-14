package sqlite

import (
	"log"
)

func GetMessages(sender string, receiver string, offset int) ([]Message, error) {
	// Initialize an empty slice of messages
	messages := []Message{}

	// Open the connection to the database.
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return messages, err
	}
	defer db.Close()

	// Query the database for all messages between the two users
	rows, err := db.Query("SELECT sender, receiver, text, timestamp FROM messages WHERE (sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?) ORDER BY message_id DESC LIMIT 10 OFFSET ?", sender, receiver, receiver, sender, offset)
	if err != nil {
		return messages, err
	}
	defer rows.Close()

	// Iterate over the rows and add each message to the slice
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.Sender, &message.Receiver, &message.Text, &message.Timestamp)
		if err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}

	// Check for any errors that may have occurred during iteration
	err = rows.Err()
	if err != nil {
		return messages, err
	}

	return messages, nil
}
