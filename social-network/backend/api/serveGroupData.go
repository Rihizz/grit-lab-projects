package api

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
)

type Group struct {
	Id          int             `json:"id"`
	CreatorID   int             `json:"creator_id"`
	CreatorName string          `json:"creator_name"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Members     []sqlite.Member `json:"members"`
	Access      bool            `json:"access"`
}

func ServeSingleGroup(w http.ResponseWriter, r *http.Request) {
	// set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	// if the request method is not GET or OPTIONS, return
	if r.Method != http.MethodGet && r.Method != http.MethodOptions {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// if the request method is OPTIONS, return
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// get the group id from the request
	groupID := r.URL.Query().Get("id")
	log.Println("groupID:", groupID)
	groupIDInt, err := strconv.Atoi(groupID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the user id from the session
	// check if the request cookie is in the sessions map
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("Error getting cookie:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session, ok := Sessions[cookie.Value]
	if !ok {
		log.Println("session not found")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	access := false
	// check if the user is a member of the group
	canViewGroup, err := sqlite.CheckIfGroupMember(session.UserID, groupIDInt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !canViewGroup {
		log.Println("User is not a member of the group")
		access = false
	} else {
		log.Println("User is a member of the group")
		access = true
	}

	// get the group data from the database
	group, err := GetGroupData(groupIDInt, session.UserID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(group)

	// get the group members from the database
	group.Members, err = getGroupMembers(groupIDInt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !access {
		log.Println("User can only see group name and description")
		// write the group data to the response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Group{
			Id:          group.Id,
			CreatorID:   group.CreatorID,
			Name:        group.Name,
			Description: group.Description,
			Access:      access,
		})

		return
	}

	// write the group data to the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Group{
		Id:          group.Id,
		CreatorID:   group.CreatorID,
		CreatorName: group.CreatorName,
		Name:        group.Name,
		Description: group.Description,
		Members:     group.Members,
		Access:      access,
	})
}

func GetGroupData(groupID, activeUser int) (Group, error) {
	canViewGroup, err := sqlite.CheckIfGroupMember(activeUser, groupID)
	if err != nil {
		log.Println(err)
		return Group{}, err
	}

	db, err := sqlite.OpenDb()
	if err != nil {
		log.Println(err)
		return Group{}, err
	}

	defer db.Close()

	if !canViewGroup {
		log.Println("User can only see group name and description")

		stmt, err := db.Prepare("SELECT id, creator_id, title, description FROM groups WHERE id = ?")
		if err != nil {
			log.Println(err)
			return Group{}, err
		}

		defer stmt.Close()

		var group Group

		err = stmt.QueryRow(groupID).Scan(&group.Id, &group.CreatorID, &group.Name, &group.Description)
		if err != nil {
			log.Println(err)
			return Group{}, err
		}

		log.Println(group)
		return group, nil
	}

	stmt, err := db.Prepare("SELECT id, creator_id, title, description FROM groups WHERE id = ?")
	if err != nil {
		log.Println(err)
		return Group{}, err
	}

	defer stmt.Close()

	var group Group
	err = stmt.QueryRow(groupID).Scan(&group.Id, &group.CreatorID, &group.Name, &group.Description)
	if err != nil {
		log.Println(err)
		return Group{}, err
	}

	log.Println(group)
	return group, nil
}

func getGroupMembers(groupID int) ([]sqlite.Member, error) {
	db, err := sqlite.OpenDb()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer db.Close()

	stmt, err := db.Prepare("SELECT user_id FROM group_members WHERE group_id = ? AND status = 1")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(groupID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	var members []sqlite.Member

	for rows.Next() {
		var member sqlite.Member
		err = rows.Scan(&member.Id)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		user, err := sqlite.GetUserById(member.Id)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		member.FullName = user.FullName

		members = append(members, member)
	}

	log.Println(members)
	return members, nil
}
