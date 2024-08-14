package pages

import (
	"groupie-tracker-visualizations/back-end/server/getData.go"
	"log"
	"net/http"
	"text/template"
)

type ArtistData struct {
	Artist         getData.Artist
	Name           string
	Image          string
	DatesLocations map[string]interface{}
	GroupMembers   []string
	CreationDate   int
	FirstAlbum     string
}

func Artist(w http.ResponseWriter, r *http.Request, index []getData.Artist) {
	data := ArtistData{}
	data.Artist = getData.GetSpecificArtist(r.URL.Path[1:], index)
	if data.Artist.Name != "" {
		data.Name = data.Artist.Name
		data.Image = data.Artist.Image
		data.DatesLocations = getData.GetRelation(data.Name, index)
		data.GroupMembers = data.Artist.Members
		data.CreationDate = data.Artist.CreationDate
		data.FirstAlbum = data.Artist.FirstAlbum
	}
	if allowedURLs(r.URL.Path, index) { // If the request is for the artist page, render the artist page
		tmpl, err := template.ParseFiles("./front-end/templates/artist.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, data)
	}
}
