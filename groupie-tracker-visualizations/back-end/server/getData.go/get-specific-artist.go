package getData

func GetSpecificArtist(artist string, arr_of_structs []Artist) Artist {
	for _, v := range arr_of_structs {
		if v.Name == artist {
			return v
		}
	}
	return Artist{}
}
