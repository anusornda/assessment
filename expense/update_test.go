package expense

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func TestUpdateExpensesById(t *testing.T) {

	body := bytes.NewBufferString(`{
		"id": 17,
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`)

	e := Expense{
		ID:     17,
		Title:  "apple smoothie",
		Amount: 89.0,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	}

	var lastest Expense

	res := request(http.MethodPut, uri("expenses", strconv.Itoa(e.ID)), body)
	err := res.Decode(&lastest)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.ID, lastest.ID)
	assert.Equal(t, e.Title, lastest.Title)
	assert.Equal(t, e.Amount, lastest.Amount)
	assert.Equal(t, e.Note, lastest.Note)
	assert.Equal(t, e.Tags, lastest.Tags)
}
