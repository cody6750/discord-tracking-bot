package functions

import (
	"net/http"
	"strings"
)

//CallAPI makes an REST API call to the given url.
func CallAPI(httpMethod, url, payload string) (*http.Response, error) {
	req, err := http.NewRequest(httpMethod, url, strings.NewReader(payload))
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
