package expense

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"net/http"
)

func (h *handler) GetExpensesHandler(c echo.Context) error {
	rows, err := h.DB.Query("SELECT id,title, amount, note, tags FROM expenses")
	if err != nil {
		return err
	}
	defer rows.Close()

	var expense []Expense
	for rows.Next() {
		var ex Expense

		err := rows.Scan(&ex.ID, &ex.Title, &ex.Amount, &ex.Note, pq.Array(&ex.Tags))

		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}

		expense = append(expense, ex)
	}

	return c.JSON(http.StatusOK, expense)
}

func (h *handler) GetExpensesByIdHandler(c echo.Context) error {

	row := h.DB.QueryRow("SELECT id,title, amount, note, tags FROM expenses WHERE id=$1", c.Param("id"))
	ex := Expense{}
	err := row.Scan(&ex.ID, &ex.Title, &ex.Amount, &ex.Note, pq.Array(&ex.Tags))

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, ex)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

}
