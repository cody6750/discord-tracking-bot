package functions

import (
	"net/http"
	"strings"
)

//CallAPI ...
func CallAPI(payload string) (*http.Response, error) {
	req, err := http.NewRequest("GET", "http://localhost:9090/crawler/item", strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err

	}
	return resp, nil
}
