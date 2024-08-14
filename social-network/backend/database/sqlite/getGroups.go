package sqlite

import "log"

type Group struct {
	Id          int
	CreatorId   int
	Title       string
	Description string
	CreatedAt   string
}

func GetGroups() ([]Group, error) {
	db, err := OpenDb()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT id, creator_id, title, description, created_at FROM groups")
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer rows.Close()

	var groups []Group

	for rows.Next() {
		var group Group
		err := rows.Scan(&group.Id, &group.CreatorId, &group.Title, &group.Description, &group.CreatedAt)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}
