package main

import (
	"log"
	"net/http"
	"social-network/api"
	"social-network/database/sqlite"
	"social-network/ws"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) //!!!!!! very nice for debugging
	// create database
	sqlite.CreateDB()

	http.HandleFunc("/api/register", api.RegisterAPI)
	http.HandleFunc("/api/login", api.LoginAPI)
	http.HandleFunc("/api/check-login", api.CheckLoginStatus)
	http.HandleFunc("/api/logout", api.LogoutAPI)
	http.HandleFunc("/api/posting", api.ServePosting)
	http.HandleFunc("/api/posts", api.ServePosts)
	http.HandleFunc("/api/user/", api.ServeUser)
	http.HandleFunc("/api/userlist", api.UserListAPI)
	http.HandleFunc("/api/create-group", api.CreateGroupAPI)
	http.HandleFunc("/api/groups", api.ServeGroups)
	http.HandleFunc("/api/user/update", api.UpdateProfileAPI)
	http.HandleFunc("/api/commenting", api.CommentingAPI)
	http.HandleFunc("/api/serve-group-data", api.ServeSingleGroup)
	http.HandleFunc("/api/serve-comments", api.ServeComments)
	http.HandleFunc("/api/group-posting", api.GroupPosting)
	http.HandleFunc("/api/serve-group-posts", api.ServeGroupPosts)
	http.HandleFunc("/api/follow", api.FollowAPI)
	http.HandleFunc("/api/post/like", api.HandlingLikes)
	http.HandleFunc("/api/event", api.CreateEventAPI)
	http.HandleFunc("/api/serve-events", api.ServeEvents)
	http.HandleFunc("/api/search-users", api.SearchUserAPI)
	http.HandleFunc("/api/serve-group-comments", api.ServeGroupComments)
	http.HandleFunc("/api/group-commenting", api.GroupCommentingAPI)
	http.HandleFunc("/api/join-group", api.JoinGroupAPI)
	http.HandleFunc("/api/event-response", api.EventResponse)
	http.HandleFunc("/api/serve-event-responses", api.ServeEventResponses)
	http.HandleFunc("/api/serve-private-messages", api.ServePrivateMessages)
	http.HandleFunc("/api/serve-group-messages", api.ServeGroupMessages)
	http.HandleFunc("/api/pposts", api.ProfilePostAPI)
	http.HandleFunc("/api/followers", api.ServeFollowerList)
	http.HandleFunc("/api/notifications", api.NotificationAPI)
	http.HandleFunc("/api/unfollow", api.UnfollowAPI)
	http.HandleFunc("/api/notif_respond", api.NotifResponse)
	http.HandleFunc("/api/followed_users", api.ServeFollowedUsers)

	ws.SetupRoutes()

	log.Println("Server running on port 6969")
	log.Fatal(http.ListenAndServe(":6969", nil))
}
