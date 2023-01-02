// / /go:build unit
// //+build unit
package expense

import (
	"bytes"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestCreateExpenses(t *testing.T) {
	// Arrange
	body := bytes.NewBufferString(`{
			"title": "strawberry smoothie",
			"amount": 79,
			"note": "night market promotion discount 10 bath",
			"tags": ["food", "beverage"]
		}`)

	t.Run("Create Expenses success", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		mockedRow := sqlmock.NewRows([]string{"id"}).AddRow(1)

		//db, mock, err := sqlmock.New()
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("INSERT INTO expenses(title, amount, note, tags) values($1, $2, $3, $4) RETURNING id").
			WithArgs("strawberry smoothie", 79.0, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"})).
			WillReturnRows(mockedRow)

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		h := handler{db}

		c := e.NewContext(req, rec)
		expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

		// Act
		err = h.CreateExpensesHandler(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("Create Expenses fail", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		mockedRow := sqlmock.NewRows([]string{"id"}).AddRow(1)

		//db, mock, err := sqlmock.New()
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("INSERT INTO expenses(title, amount, note, tags) values($1, $2, $3, $4) RETURNING id").
			WithArgs("strawberry smoothie", 79.0, "night market promotion discount 10 bath", "[\"food\", \"beverage\"]").
			WillReturnRows(mockedRow)

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		h := handler{db}

		c := e.NewContext(req, rec)

		// Act
		err = h.CreateExpensesHandler(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

}

func TestUpdateExpensesById(t *testing.T) {

	body := bytes.NewBufferString(`{
			"title": "apple smoothie",
			"amount": 89,
			"note": "no discount",
			"tags": ["beverage"]
		}`)

	bodyBadRequest := bytes.NewBufferString(`{
			"title": "apple smoothie",
			"amount": "89"",
			"note": "no discount",
			"tags": ["beverage"]
		}`)

	t.Run("update Expenses success", func(t *testing.T) {
		id := 1
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/", body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		result := sqlmock.NewResult(1, 1)

		//db, mock, err := sqlmock.New()
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectPrepare("UPDATE expenses SET title=$2 , amount=$3, note=$4, tags=$5 WHERE id=$1;").
			ExpectExec().
			WithArgs(id, "apple smoothie", 89.0, "no discount", pq.Array([]string{"beverage"})).
			WillReturnResult(result)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		h := handler{db}

		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(id))
		expected := "{\"id\":1,\"title\":\"apple smoothie\",\"amount\":89,\"note\":\"no discount\",\"tags\":[\"beverage\"]}"

		// Act
		err = h.UpdateExpensesByIdHandler(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("update Expenses bad request", func(t *testing.T) {
		id := 1
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/", bodyBadRequest)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		result := sqlmock.NewResult(1, 1)

		//db, mock, err := sqlmock.New()
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectPrepare("UPDATE expenses SET title=$2 , amount=$3, note=$4, tags=$5 WHERE id=$1;").
			ExpectExec().
			WithArgs("id", "apple smoothie", 89.0, "no discount", pq.Array([]string{"beverage"})).
			WillReturnResult(result)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		h := handler{db}

		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(id))

		// Act
		err = h.UpdateExpensesByIdHandler(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("update Expenses fail", func(t *testing.T) {
		id := 1
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/", bodyBadRequest)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		result := sqlmock.NewResult(1, 1)

		//db, mock, err := sqlmock.New()
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectPrepare("UPDATE expenses SET title=$2 , amount=$3, note=$4, tags=$5 WHERE id=$1;").
			ExpectExec().
			WithArgs("id", "apple smoothie", "no discount", pq.Array([]string{"beverage"})).
			WillReturnResult(result)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		h := handler{db}

		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(id))

		// Act
		err = h.UpdateExpensesByIdHandler(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

}

func TestGetAllExpenses(t *testing.T) {

	t.Run("get all expense", func(t *testing.T) {
		// Arrange

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		newsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow(1, "apple smoothie", 89, "no discount", pq.Array([]string{"beverage"})).
			AddRow(2, "iPhone 14 Pro Max 1TB", 66900, "birthday gift from my love", pq.Array([]string{"gadget"}))

		db, mock, err := sqlmock.New()
		mock.ExpectQuery("SELECT id,title, amount, note, tags FROM expenses").WillReturnRows(newsMockRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		h := handler{db}
		c := e.NewContext(req, rec)
		expected := "[{\"id\":1,\"title\":\"apple smoothie\",\"amount\":89,\"note\":\"no discount\",\"tags\":[\"beverage\"]}," +
			"{\"id\":2,\"title\":\"iPhone 14 Pro Max 1TB\",\"amount\":66900,\"note\":\"birthday gift from my love\",\"tags\":[\"gadget\"]}]"

		// Act
		err = h.GetExpensesHandler(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("get expense by id", func(t *testing.T) {
		id := 1
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		newsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow(1, "strawberry smoothie", 79, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"}))

		//db, mock, err := sqlmock.New()
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("SELECT id,title, amount, note, tags FROM expenses WHERE id=$1").
			WithArgs(strconv.Itoa(id)).
			WillReturnRows(newsMockRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		h := handler{db}

		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(id))
		expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

		// Act
		err = h.GetExpensesByIdHandler(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}

	})
}
