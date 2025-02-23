//go:build integration
// +build integration

package expense

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type Response struct {
	*http.Response
	err error
}

var serverPort = os.Getenv("PORT")
var dbUrl = os.Getenv("DATABASE_URL")

func TestITCreateExpenses(t *testing.T) {
	// Arrange
	body := bytes.NewBufferString(`{
				"title": "strawberry smoothie",
				"amount": 79,
				"note": "night market promotion discount 10 bath",
				"tags": ["food", "beverage"]
			}`)

	t.Run("Create Expenses success", func(t *testing.T) {
		h := InitDB(dbUrl)
		eh := initialEcho()

		eh.POST("expenses", h.CreateExpensesHandler)

		go func(e *echo.Echo) {
			eh.Start(fmt.Sprintf("%v", serverPort))
		}(eh)

		for {
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost%v", serverPort), 30*time.Second)
			if err != nil {
				log.Println(err)
			}
			if conn != nil {
				conn.Close()
				break
			}
		}

		// Arrange
		var ex Expense

		res := request(http.MethodPost, uri("expenses"), body)
		err := res.Decode(&ex)

		// Assertions

		expected := Expense{
			Title:  "strawberry smoothie",
			Amount: 79.0,
			Note:   "night market promotion discount 10 bath",
			Tags:   []string{"food", "beverage"},
		}

		if assert.NoError(t, err) {
			assert.EqualValues(t, http.StatusCreated, res.StatusCode)
			assert.Equal(t, expected.Title, ex.Title)
			assert.Equal(t, expected.Amount, ex.Amount)
			assert.Equal(t, expected.Note, ex.Note)
			assert.Equal(t, expected.Tags, ex.Tags)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = eh.Shutdown(ctx)
		assert.NoError(t, err)

	})

}

func TestITGetExpensesById(t *testing.T) {

	h := InitDB(dbUrl)
	eh := initialEcho()

	eh.POST("expenses", h.CreateExpensesHandler)
	eh.GET("/expenses/:id", h.GetExpensesByIdHandler)

	go func(e *echo.Echo) {
		eh.Start(fmt.Sprintf("%v", serverPort))
	}(eh)

	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost%v", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	// Arrange
	ex := seedExpenses(t)
	var lastest Expense

	res := request(http.MethodGet, uri("expenses", strconv.Itoa(ex.ID)), nil)
	err := res.Decode(&lastest)

	// Assertions

	if assert.NoError(t, err) {
		assert.EqualValues(t, http.StatusOK, res.StatusCode)
		assert.EqualValues(t, ex.ID, lastest.ID)
		assert.Equal(t, ex.Title, lastest.Title)
		assert.Equal(t, ex.Amount, lastest.Amount)
		assert.Equal(t, ex.Note, lastest.Note)
		assert.Equal(t, ex.Tags, lastest.Tags)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestUpdateExpensesById(t *testing.T) {

	h := InitDB(dbUrl)
	eh := initialEcho()

	eh.POST("expenses", h.CreateExpensesHandler)
	eh.PUT("/expenses/:id", h.UpdateExpensesByIdHandler)

	go func(e *echo.Echo) {
		eh.Start(fmt.Sprintf("%v", serverPort))
	}(eh)

	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost%v", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	// Arrange
	ex := seedExpenses(t)

	body := bytes.NewBufferString(`{
				"title": "apple smoothie",
				"amount": 89,
				"note": "no discount",
				"tags": ["beverage"]
			}`)

	var lastest Expense

	res := request(http.MethodPut, uri("expenses", strconv.Itoa(ex.ID)), body)
	err := res.Decode(&lastest)

	// Assertions

	expected := Expense{
		ID:     ex.ID,
		Title:  "apple smoothie",
		Amount: 89.0,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	}

	if assert.NoError(t, err) {
		assert.EqualValues(t, http.StatusOK, res.StatusCode)
		assert.EqualValues(t, expected.ID, lastest.ID)
		assert.Equal(t, expected.Title, lastest.Title)
		assert.EqualValues(t, expected.Amount, lastest.Amount)
		assert.Equal(t, expected.Note, lastest.Note)
		assert.Equal(t, expected.Tags, lastest.Tags)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestITGetAllExpenses(t *testing.T) {

	h := InitDB(dbUrl)
	eh := initialEcho()

	eh.GET("/expenses", h.GetExpensesHandler)

	go func(e *echo.Echo) {
		eh.Start(fmt.Sprintf("%v", serverPort))
	}(eh)

	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost%v", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	// Arrange
	var ex []Expense
	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&ex)

	// Assertions

	if assert.NoError(t, err) {
		assert.EqualValues(t, http.StatusOK, res.StatusCode)
		assert.Greater(t, len(ex), 0)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestGetAllExpensesUnauthorized(t *testing.T) {

	h := InitDB(dbUrl)
	eh := initialEcho()

	eh.GET("/expenses", h.GetExpensesHandler)

	go func(e *echo.Echo) {
		eh.Start(fmt.Sprintf("%v", serverPort))
	}(eh)

	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost%v", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

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

	err := request(http.MethodPost, fmt.Sprintf("http://localhost%v/expenses", serverPort), body).Decode(&ex)
	if err != nil {
		t.Fatal("can't create expense:", err)
	}
	return ex
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization", "November 10, 2009")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

func requestUnauthorized(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization", "November 10, 2009wrong_token")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

func uri(paths ...string) string {
	host := fmt.Sprintf("http://localhost%v", serverPort)

	if paths == nil {
		return host
	}
	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func initialEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(CheckUserAuth())

	return e
}
