package expense

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	*http.Response
	err error
}

func uri(paths ...string) string {
	host := "http://localhost:2565"

	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "November 10, 2009")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

func requestUnauthorized(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "November 10, 2009wrong_token")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}
