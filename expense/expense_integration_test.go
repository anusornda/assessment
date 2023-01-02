//go:build integration
// +build integration

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

func TestUpdateExpensesById(t *testing.T) {

	e := seedExpenses(t)

	body := bytes.NewBufferString(`{
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`)

	var lastest Expense

	res := request(http.MethodPut, uri("expenses", strconv.Itoa(e.ID)), body)
	err := res.Decode(&lastest)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.ID, lastest.ID)
	assert.Equal(t, "apple smoothie", lastest.Title)
	assert.Equal(t, 89.0, lastest.Amount)
	assert.Equal(t, "no discount", lastest.Note)
	assert.Equal(t, []string{"beverage"}, lastest.Tags)
}

func TestITGetExpensesById(t *testing.T) {
	e := seedExpenses(t)

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

func TestITGetAllExpenses(t *testing.T) {

	var ex []Expense
	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&ex)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(ex), 0)
}

func TestGetAllExpensesUnauthorized(t *testing.T) {

	var ex []Expense
	res := requestUnauthorized(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&ex)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, res.StatusCode)
}

func seedExpenses(t *testing.T) Expense {
	var ex Expense
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath",
		"tags": ["food", "beverage"]
	}`)

	err := request(http.MethodPost, uri("expenses"), body).Decode(&ex)
	if err != nil {
		t.Fatal("can't create expense:", err)
	}
	return ex
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
