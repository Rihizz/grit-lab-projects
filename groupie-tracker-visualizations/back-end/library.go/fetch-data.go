package library

import (
	"io/ioutil"
	"net/http"
)

// library.FetchData fetches the data from a given URL with json and returns it as an array of bytes
// and an error
func FetchData(url string) ([]byte, error) {
	response, err := http.Get(url)

	if err != nil {
		return []byte{}, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	return responseData, err
}
