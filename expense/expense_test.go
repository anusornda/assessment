package expense

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

type Response struct {
	*http.Response
	err error
}

func TestCreateExpenses(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)

	var ex Expense

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&ex)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "strawberry smoothie", ex.Title)
	assert.Equal(t, 79.0, ex.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", ex.Note)
	assert.Equal(t, []string{"food", "beverage"}, ex.Tags)
}

func TestGetExpensesById(t *testing.T) {
	e := Expense{
		ID:     17,
		Title:  "strawberry smoothie",
		Amount: 79.0,
		Note:   "night market promotion discount 10 bath",
		Tags:   []string{"food", "beverage"},
	}

	var lastest Expense

	res := request(http.MethodGet, uri("expenses", strconv.Itoa(e.ID)), nil)
	err := res.Decode(&lastest)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.Title, lastest.Title)
	assert.Equal(t, e.Amount, lastest.Amount)
	assert.Equal(t, e.Note, lastest.Note)
	assert.Equal(t, e.Tags, lastest.Tags)
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
	//req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}
