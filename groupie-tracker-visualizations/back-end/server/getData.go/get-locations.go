package getData

import (
	"encoding/json"
	"groupie-tracker-visualizations/back-end/library.go"
)

type Locations struct {
	Id     int      `json:"id"`
	Places []string `json:"locations"`
	Dates  string   `json:"dates"`
}

func GetLocation(artist string, arr_of_structs []Artist) Locations {
	locations := Locations{}
	for _, v := range arr_of_structs {
		if v.Name == artist {
			data, _ := library.FetchData(v.Locations)
			json.Unmarshal(data, &locations)
			break
		}
	}
	return locations
}
