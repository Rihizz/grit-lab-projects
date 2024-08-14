package sqlite

import "log"

func AcceptGroup(groupID, sender int) error {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	_, err = db.Exec("UPDATE group_members SET status = 1 WHERE user_id = ? AND group_id", sender, groupID)
	if err != nil {
    log.Fatal(err)
	}

	return nil
}