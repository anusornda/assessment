package expense

import (
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"net/http"
)

func CreateExpensesHandler(c echo.Context) error {

	var ex Expense
	err := c.Bind(&ex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := db.QueryRow("INSERT INTO expenses(title, amount, note, tags) values($1, $2, $3, $4) RETURNING id", ex.Title, ex.Amount, ex.Note, pq.Array(ex.Tags))
	err = row.Scan(&ex.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, ex)

}
