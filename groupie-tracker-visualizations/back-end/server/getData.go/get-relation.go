package getData

import (
	"encoding/json"
	"groupie-tracker-visualizations/back-end/library.go"
)

type Relation struct {
	Index []struct {
		Id             int                    `json:"id"`
		DatesLocations map[string]interface{} `json:"datesLocations"`
	} `json:"index"`
}

func GetRelation(artist string, arr_of_structs []Artist) map[string]interface{} {
	data, _ := library.FetchData("https://groupietrackers.herokuapp.com/api/relation")
	arr := Relation{}
	json.Unmarshal(data, &arr)
	for _, v := range arr_of_structs {
		if v.Name == artist {
			old := arr.Index[v.Id-1].DatesLocations
			result := map[string]interface{}{}
			for k, v := range old {
				result[library.TextBeautifier(k)] = v
			}
			return result
		}
	}
	return map[string]interface{}{}
}
