package ws

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"social-network/api"
	"social-network/database/sqlite"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type UserConnection struct {
	UserID     int
	Username   string
	Connection *websocket.Conn
}

type ChatConnection struct {
	UserID	 int
	Username string
	Connection *websocket.Conn
}

var ChatConnections = make(map[*websocket.Conn]*ChatConnection)
var ChatConnectionsByName = make(map[string]*websocket.Conn)

var Connections = make(map[*websocket.Conn]*UserConnection)
var ConnectionsByName = make(map[string]*websocket.Conn)

type Message struct {
	Command    string `json:"command"`
	Text       string `json:"message"`
	Receiver   string `json:"receiver"`
	ReceiverID int    `json:"receiver_id"`
	Sender     string `json:"sender"`
	SenderID   int    `json:"sender_id"`
	Image      string `json:"image"`
	Timestamp  string
}

type Notification struct {
	Type      string `json:"Type"`
	Groupid   int    `json:"Groupid"`
	UserID    int    `json:"UserID"`
	SenderID  int    `json:"Sender"`
	Message   string `json:"Message"`
	Command   string `json:"Command"`
	Timestamp string `json:"Date"`
}

type NotificationResponse struct {
	NotifId  int    `json:"NotifId"`
	UserId   int    `json:"UserId"`
	SenderId int    `json:"SenderId"`
	Groupid  int    `json:"Groupid"`
	Content  string `json:"Content"`
	Read     bool   `json:"Read"`
	Type     string `json:"Type"`
	Date     string `json:"Date"`
	Command  string `json:"Command"`
}

// function to read the data from the websocket connection
func reader(conn *websocket.Conn) {
	// Set up a close handler for the WebSocket connection
	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("WebSocket closed with code %d and text: %s", code, text)
		delete(Connections, conn) // Remove the connection from the map.
		return nil
	})
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// Create a map to store the connections by username
		for _, conn := range Connections {
			ConnectionsByName[conn.Username] = conn.Connection
		}
		log.Println("message received: ")
		log.Println("messageType: ", messageType)
		log.Println(string(p))
		var message Message
		err = json.Unmarshal(p, &message)
		if err != nil {
			log.Println(err)
			return
		}
		// if message.Command == "NEW_MESSAGE" {
		// 	handleNewMessage(conn, message)
		// }
		// if message.Command == "GROUP_MESSAGE" {
		// 	handleGroupMessage(conn, message)
		// }
		if message.Command == "NOTIFICATION" {
			handleNotification(conn, message)
		}
	}
}

func chatReader(conn *websocket.Conn) {
	// Set up a close handler for the WebSocket connection
	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("WebSocket closed with code %d and text: %s", code, text)
		delete(ChatConnections, conn) // Remove the connection from the map.
		return nil
	})
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// Create a map to store the connections by username
		for _, conn := range ChatConnections {
			ChatConnectionsByName[conn.Username] = conn.Connection
		}
		log.Println("message received: ")
		log.Println("messageType: ", messageType)
		log.Println(string(p))
		var message Message
		err = json.Unmarshal(p, &message)
		if err != nil {
			log.Println(err)
			return
		}
		if message.Command == "NEW_MESSAGE" {
			handleNewMessage(conn, message)
		}
		if message.Command == "GROUP_MESSAGE" {
			handleGroupMessage(conn, message)
		}
	}
}

// function to set up the websocket endpoint
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	c, err := r.Cookie("session_token")
	if err != nil {
		log.Println("User is not logged in")
		return
	}

	userInfo, ok := api.Sessions[c.Value]
	if !ok {
		log.Println("Session error user is not logged in")
		return
	}

	userconn := &UserConnection{
		UserID:     userInfo.UserID,
		Username:   userInfo.Username,
		Connection: ws,
	}

	// Add the connection to the list of active connections.
	Connections[ws] = userconn

	log.Printf("User %s with ID %d successfully connected", userconn.Username, userconn.UserID)
	log.Println("connections: ", Connections)
	go reader(ws)
}

// function to handle the chat endpoint
func chatEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	c, err := r.Cookie("session_token")
	if err != nil {
		log.Println("User is not logged in")
		return
	}

	userInfo, ok := api.Sessions[c.Value]
	if !ok {
		log.Println("Session error user is not logged in")
		return
	}

	userconn := &ChatConnection{
		UserID:     userInfo.UserID,
		Username:   userInfo.Username,
		Connection: ws,
	}

	// Add the connection to the list of active connections.
	ChatConnections[ws] = userconn
	// for c := range Connections {
	// 	if c != ws {
	// 		err := c.WriteMessage(websocket.TextMessage, []byte("USER_JOINED"))
	// 		if err != nil {
	// 			log.Println(err)
	// 			delete(Connections, c) // Remove the connection from the map.
	// 		}
	// 	}
	// }
	log.Printf("User %s with ID %d successfully connected", userconn.Username, userconn.UserID)
	log.Println("connections: ", ChatConnections)
	go chatReader(ws)
}


// function to set up the routes for the websocket server.
func SetupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/chat", chatEndpoint)
}

func handleNewMessage(conn *websocket.Conn, message Message) {
	log.Printf("Received message: %v\n", message)
	// get the sender username
	sender, ok := ChatConnections[conn]
	if !ok {
		log.Println("Sender not found")
	}

	senderUsername := sender.Username
	message.Sender = senderUsername
	log.Println("Sender: ", senderUsername)
	message.Timestamp = time.Now().Format("2006-01-02 15:04:05")

	// Decode the image data from Base64
	var imageData []byte
	if message.Image != "" {
		splitImage := strings.Split(message.Image, ",")
		message.Image = splitImage[1]
		var err error
		imageData, err = base64.StdEncoding.DecodeString(message.Image)
		if err != nil {
			log.Println("Error decoding image data: ", err)
			return
		}
	}

	// add message to database
	sqlite.AddPrivateMessage(message.SenderID, message.ReceiverID, message.Text, imageData, message.Timestamp)
	log.Println("Message: ", message)
	// Send message to sender
	log.Println("Sending message to sender: ", message)
	err := conn.WriteJSON(message)
	if err != nil {
		log.Println(err)
		return
	}
	// Check if receiver is online
	log.Println("ChatConnectionsByName: ", ChatConnectionsByName)
	for c := range ConnectionsByName {
		log.Println(c)
	}
	receiverConn, ok := ChatConnectionsByName[message.Receiver]
	if !ok {
		log.Printf("Receiver %v is not online, message will be saved to database\n", message.Receiver)
		return
	}
	// Send message to receiver
	if message.Receiver != message.Sender {
		log.Println("Sending message to receiver: ", message)
		err = receiverConn.WriteJSON(message)
		if err != nil {
			log.Println(err)
		}
	}
}

func handleGroupMessage(conn *websocket.Conn, message Message) {
	sender, ok := ChatConnections[conn]
	if !ok {
		log.Println("Sender not found")
		return
	}

	senderUsername := sender.Username
	message.Sender = senderUsername
	log.Println("Sender: ", senderUsername)
	message.Timestamp = time.Now().Format("2006-01-02 15:04:05")

	var imageData []byte
	if message.Image != "" {
		splitImage := strings.Split(message.Image, ",")
		message.Image = splitImage[1]
		var err error
		imageData, err = base64.StdEncoding.DecodeString(message.Image)
		if err != nil {
			log.Println("Error decoding image data: ", err)
			return
		}
	}

	// add message to database
	groupID, err := strconv.Atoi(message.Receiver)
	if err != nil {
		log.Println(err)
		return
	}

	sqlite.AddGroupMessage(groupID, message.SenderID, message.Text, imageData, message.Timestamp)

	log.Println("Message: ", message)
	// Send message to sender
	log.Println("Sending message to sender: ", message)
	err = conn.WriteJSON(message)
	if err != nil {
		log.Println(err)
		return
	}

	// get all users in the group
	members, err := sqlite.GetGroupMembers(groupID)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("members: ", members)

	// Send message to all group members
	for _, member := range members {
		if member.FullName != message.Sender {
			receiverConn, ok := ChatConnectionsByName[member.FullName]
			if !ok {
				log.Printf("Receiver %v is not online, message will be saved to database\n", message.Receiver)
				continue
			}
			log.Println("Sending message to receiver: ", message)
			err = receiverConn.WriteJSON(message)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func handleNotification(conn *websocket.Conn, message Message) {
	// Unmarshal the JSON into a Notification struct

	log.Println("Received notification: ", message)
	var notification Notification
	err := json.Unmarshal([]byte(message.Text), &notification)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Notification: ", notification)

	receiver, err := sqlite.GetUserById(notification.UserID)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("sender vitun id", notification.SenderID)

	sender, err := sqlite.GetUserById(notification.SenderID)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Sender: ", sender)

	log.Println("group id: ", notification.Groupid)

	if notification.SenderID == notification.UserID {
		log.Println("Sender and receiver are the same, notification will not be sent")
		return
	}

	xd, err := checkNotifications(notification.UserID, notification.SenderID, notification.Message)
	if err != nil {
		log.Println(err)
		return
	}
	if xd {
		log.Println("duplicate notification, will not be sent")
		return
	}

	// add notification to database
	id, err := addNotification(notification.UserID, notification.SenderID, notification.Groupid, notification.Message, notification.Type)
	if err != nil {
		log.Println(err)
		return
	}

	notification.Command = "NOTIFICATION"

	receiverConn, ok := ConnectionsByName[receiver.FullName]
	if !ok {
		log.Printf("Receiver %v is not online, notification will be saved to database\n", receiver.FullName)
		return
	}

	var notifResponse NotificationResponse
	notifResponse.UserId = notification.UserID
	notifResponse.SenderId = notification.SenderID
	notifResponse.Content = notification.Message
	notifResponse.Type = notification.Type
	notifResponse.Date = time.Now().Format("2006-01-02 15:04:05")
	notifResponse.Read = false
	notifResponse.NotifId = id
	notifResponse.Command = "NOTIFICATION"
	notifResponse.Groupid = notification.Groupid

	err = receiverConn.WriteJSON(notifResponse)
	if err != nil {
		log.Println(err)
	}

}

func addNotification(userID int, SenderID int, Groupid int, content string, nType string) (int, error) {
	// open database connection
	db, err := sqlite.OpenDb()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	// close database connection
	defer db.Close()

	// insert notification into database
	var result sql.Result
	if Groupid != 0 {
		result, err = db.Exec("INSERT INTO notifications (user_id, sender_id, reference_id, content, type) VALUES (?, ?, ?, ?, ?)", userID, SenderID, Groupid, content, nType)
	} else {
		result, err = db.Exec("INSERT INTO notifications (user_id, sender_id, content, type) VALUES (?, ?, ?, ?)", userID, SenderID, content, nType)
	}
	if err != nil {
		log.Println(err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return int(id), nil
}

func checkNotifications(UserId int, SenderId int, message string) (bool, error) {
	// open database connection
	db, err := sqlite.OpenDb()
	if err != nil {
		log.Println(err)
		return true, err
	}

	// close database connection
	defer db.Close()

	// check if notification already exists with the same message
	var id int
	err = db.QueryRow("SELECT id FROM notifications WHERE user_id = ? AND sender_id = ? AND content = ?", UserId, SenderId, message).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Notification does not exist")
			return false, nil
		}
		log.Println(err)
		return true, err
	}
	return true, nil
}
