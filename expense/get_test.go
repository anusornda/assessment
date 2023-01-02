package expense

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func TestGetExpensesById(t *testing.T) {
	e := Expense{
		ID:     15,
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

func TestGetAllExpenses(t *testing.T) {

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
