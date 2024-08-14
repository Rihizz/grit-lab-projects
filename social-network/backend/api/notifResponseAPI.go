package api

import (
	"fmt"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
)

func NotifResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != "POST" && r.Method != "OPTIONS" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "{\"status\": 405, \"message\": \"method not allowed\"}")
		return
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"status\": 200, \"message\": \"success\"}")
		return
	}

	// get the user id from the session
	//check if the request cookie is in the sessions map
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
	userID := session.UserID
	log.Println("user id:", userID)

	//get the notification id from the url
	notificationID := r.URL.Query().Get("notif_id")
	log.Println("notification id:", notificationID)

	notifID, err := strconv.Atoi(notificationID)
	if err != nil {
		log.Println("Error converting notification id to int:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//get the notification type from the url
	notificationType := r.URL.Query().Get("notif_type")
	log.Println("notification type:", notificationType)

	SenderID := r.URL.Query().Get("sender_id")
	log.Println("sender id:", SenderID)

	sender, err := strconv.Atoi(SenderID)
	if err != nil {
		log.Println("Error converting sender id to int:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the notification response from the url
	notificationResponse := r.URL.Query().Get("response")
	log.Println("notification response:", notificationResponse)

	if notificationType == "follow request" {
		log.Println("follow request")
		if notificationResponse == "accept" {
			sqlite.AcceptFollow(userID, sender)
			sqlite.DeleteNotification(notifID)
			return
		} else if notificationResponse == "decline" {
			sqlite.DeclineFollow(userID, sender)
			sqlite.DeleteNotification(notifID)
			return
		}
	}

	group := r.URL.Query().Get("group_id")
	groupId, err := strconv.Atoi(group)
	if err != nil {
		log.Println("Error converting group id to int:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("group id:", groupId)

	if notificationType == "group request" {
		log.Println("group request")
		if notificationResponse == "accept" {
			sqlite.AcceptGroup(groupId, sender)
			sqlite.DeleteNotification(notifID)
			return
		} else if notificationResponse == "decline" {
			sqlite.DeclineGroup(groupId, sender)
			sqlite.DeleteNotification(notifID)
			return
		}
	}

	if notificationType == "group invite" {
		log.Println("group invite")
		if notificationResponse == "accept" {
			sqlite.AddGroupMember(userID, groupId, 1)
			sqlite.DeleteNotification(notifID)
			return
		} else if notificationResponse == "decline" {
			sqlite.DeleteNotification(notifID)
			return
		}
	}

	// delete the notification from the database
	sqlite.DeleteNotification(notifID)

}
