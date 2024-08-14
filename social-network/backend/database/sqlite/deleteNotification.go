package sqlite

import "log"

func DeleteNotification(notifId int) error {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM notifications WHERE id = ?", notifId)
	if err != nil {
		log.Println(err)
	}
	return nil		
}