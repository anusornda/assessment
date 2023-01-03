//go:build integration
// +build integration

package expense

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

type Response struct {
	*http.Response
	err error
}

const serverPort = 80

func TestITCreateExpenses(t *testing.T) {
	// Arrange
	body := bytes.NewBufferString(`{
				"title": "strawberry smoothie",
				"amount": 79,
				"note": "night market promotion discount 10 bath",
				"tags": ["food", "beverage"]
			}`)

	t.Run("Create Expenses success", func(t *testing.T) {
		eh := echo.New()
		go func(e *echo.Echo) {
			db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
			if err != nil {
				log.Fatal(err)
			}

			h := NewHandler(db)

			e.POST("expenses", h.CreateExpensesHandler)
			e.Start(fmt.Sprintf(":%d", serverPort))
		}(eh)
		for {
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
			if err != nil {
				log.Println(err)
			}
			if conn != nil {
				conn.Close()
				break
			}
		}

		//req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), body)
		//assert.NoError(t, err)
		//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		//client := http.Client{}
		//
		//// Act
		//resp, err := client.Do(req)
		//assert.NoError(t, err)
		//
		//byteBody, err := ioutil.ReadAll(resp.Body)
		//assert.NoError(t, err)
		//resp.Body.Close()
		//
		//// Assertions
		//expected := "[{\"ID\":1,\"Title\":\"test-title\",\"Content\":\"test-content\",\"Author\":\"test-author\"}]"
		//
		//if assert.NoError(t, err) {
		//	assert.Equal(t, http.StatusCreated, resp.StatusCode)
		//	assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
		//}

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
	//host := "http://localhost:2565"
	host := fmt.Sprintf("http://localhost:%d", serverPort)

	if paths == nil {
		return host
	}
	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}
