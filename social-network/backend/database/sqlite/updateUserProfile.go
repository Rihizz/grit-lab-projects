package sqlite

import (
	"bytes"
	"fmt"
	"log"
)

func UpdateUserProfile(userID, email, nickname, aboutMe, fileName string, fileContent []byte, newPassword string, privacy string) error {
	db, err := OpenDb()
	if err != nil {
		return err
	}
	defer db.Close()

	// Update user data
	var updateFields bytes.Buffer
	var updateValues []interface{}

	if newPassword != "" {
		password, err := hashPassword(newPassword)
		if err != nil {
			log.Println(err)
			//tx.Rollback() is this essential?
			return err
		}

		updateFields.WriteString(" password = ?,")
		updateValues = append(updateValues, password)
	}

	if email != "" {
		updateFields.WriteString("email = ?,")
		updateValues = append(updateValues, email)
	}

	if nickname != "" {
		updateFields.WriteString(" nickname = ?,")
		updateValues = append(updateValues, nickname)
	}

	if aboutMe != "" {
		updateFields.WriteString(" aboutme = ?,")
		updateValues = append(updateValues, aboutMe)
	}

	if privacy != "" {
		updateFields.WriteString(" privacy = ?,")
		updateValues = append(updateValues, privacy)
	}

	// Prepare the statement only if there are fields to update
	if updateFields.Len() > 0 {

		// Remove the trailing comma
		updateFields.Truncate(updateFields.Len() - 1)
		stmt, err := db.Prepare(fmt.Sprintf("UPDATE users SET %s WHERE user_id = ?", updateFields.String()))
		if err != nil {
			return err
		}
		updateValues = append(updateValues, userID)
		_, err = stmt.Exec(updateValues...)
		if err != nil {
			return err
		}
	}

	// Update avatar and avatarname if a new one was uploaded
	if fileName != "" {
		stmt, err := db.Prepare("UPDATE users SET avatar = ?, avatarname = ? WHERE user_id = ?")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(fileContent, fileName, userID)
		if err != nil {
			return err
		}
	}
	return nil
}
