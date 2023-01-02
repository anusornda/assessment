package expense

import (
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"net/http"
)

func (h *handler) CreateExpensesHandler(c echo.Context) error {

	var ex Expense
	err := c.Bind(&ex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := h.DB.QueryRow("INSERT INTO expenses(title, amount, note, tags) values($1, $2, $3, $4) RETURNING id", ex.Title, ex.Amount, ex.Note, pq.Array(ex.Tags))
	err = row.Scan(&ex.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, ex)

}

func (h *handler) Greeting(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
