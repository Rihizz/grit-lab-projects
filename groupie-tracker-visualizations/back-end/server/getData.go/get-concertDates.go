package getData

import (
	"encoding/json"
	"groupie-tracker-visualizations/back-end/library.go"
)

type ConcertDates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

func GetConcertDates(artist string, arr_of_structs []Artist) []string {
	result := ConcertDates{}
	for _, v := range arr_of_structs {
		if v.Name == artist {
			data, _ := library.FetchData(v.ConcertDates)
			json.Unmarshal(data, &result)
			break
		}
	}
	return result.Dates
}
