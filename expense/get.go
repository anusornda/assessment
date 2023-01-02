package expense

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"net/http"
)

func GetExpensesHandler(c echo.Context) error {

	stmt, err := db.Prepare("SELECT id,title, amount, note, tags FROM expenses")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	rows, err := stmt.Query()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	expense := []Expense{}
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

func GetExpensesByIdHandler(c echo.Context) error {

	stmt, err := db.Prepare("SELECT id,title, amount, note, tags FROM expenses WHERE  id=$1")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	row := stmt.QueryRow(c.Param("id"))
	ex := Expense{}
	err = row.Scan(&ex.ID, &ex.Title, &ex.Amount, &ex.Note, pq.Array(&ex.Tags))

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, ex)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

}
