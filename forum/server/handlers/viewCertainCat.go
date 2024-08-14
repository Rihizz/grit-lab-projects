package handlers

import (
	"forum/server/login"
	"forum/server/sqlite"
	util "forum/server/utils"
	"log"
	"net/http"
	"sort"
	"strings"
	"text/template"
)

func viewCertainCat(w http.ResponseWriter, r *http.Request) {
	// Check if the URL path is allowed
	if allowedBoards(r.URL.Path, sqlite.GetCategories()) {
		// load the template
		tmpl, err := template.ParseFiles("templates/posts.html")
		if err != nil {
			log.Println("template parsing error: ", err)
			ErrorHandling(500, w, r)
			return
		}
		catposts := certainCat(r.URL.Path[1:])
		// Check if the user is logged in
		if login.CheckSession(w, r) {
			c, err := r.Cookie("session_token")
			if err != nil {
				log.Println("viewCertainCat - 500. Can't get cookie.")
				ErrorHandling(500, w, r)
				return
			}
			userd := login.Sessions[c.Value]
			tmpl.Execute(w, map[string]interface{}{
				"posts":    catposts,
				"category": r.URL.Path[1:],
				"user":     userd,
			})
		} else {
			tmpl.Execute(w, map[string]interface{}{
				"posts":    catposts,
				"category": r.URL.Path[1:],
			})
		}
	} else {
		// if the URL path is not allowed, return 404
		ErrorHandling(404, w, r)
	}
}

func certainCat(cat string) []util.Fpost {
	var matchingPosts []util.Fpost
	for _, post := range util.Fposts {
		// Split the Category field by space
		categories := strings.Split(post.Category, " ")

		// Check if the category is in the list of categories
		if contains(categories, cat) {
			// Add the post to the matchingPosts slice
			matchingPosts = append(matchingPosts, post)
		}
	}
	// Sort the posts by ID
	sort.Slice(matchingPosts, func(i, j int) bool {
		return matchingPosts[i].PostID < matchingPosts[j].PostID
	})
	return matchingPosts
}

// contains checks if the slice contains the item
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
