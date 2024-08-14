package ws

// import (
// 	"log"

// 	"github.com/gorilla/websocket"
// )

// type JoinGroupRequest struct {
// 	GroupId   string `json:"group_id"`
// 	CreatorID string `json:"creator_id"`
// 	UserID    string `json:"user_id"`
// }

// // type Notification struct {
// // 	Command   string `json:"command"`
// // 	GroupId   string `json:"group_id"`
// // 	GroupName string `json:"group_name"`
// // 	UserId    string `json:"user_id"`
// // 	Username  string `json:"username"`
// // }

// func JoinGroupRequestHandler(conn *websocket.Conn) {
// 	log.Println("JoinGroupRequest")
// 	sender, ok := Connections[conn]
// 	if !ok {
// 		log.Println("Connection not found")
// 		return
// 	}

// 	var request JoinGroupRequest
// 	err := conn.ReadJSON(&request)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	// get the group id from the request
// 	groupID := request.GroupId

// 	// send a notification to the group owner
// 	// get the group owner's connection
// 	ownerConn, ok := ConnectionsByName[request.CreatorID]
// 	if !ok {
// 		log.Println("Connection not found")
// 		return
// 	}

// 	// send the notification
// 	err = ownerConn.WriteJSON(Notification{
// 		Command:   "JOIN_GROUP_REQUEST",
// 		GroupId:   groupID,
// 		GroupName: request.GroupId,
// 		UserId:    request.UserID,
// 		Username:  sender.Username,
// 	})
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// }
