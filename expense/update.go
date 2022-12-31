package expense

import (
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"net/http"
	"strconv"
)

func UpdateExpensesByIdHandler(c echo.Context) error {

	rowID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	ex := Expense{}
	err = c.Bind(&ex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	stmt, err := db.Prepare("UPDATE expenses SET title=$2 , amount=$3, note=$4, tags=$5 WHERE id=$1;")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	if _, err := stmt.Exec(rowID, ex.Title, ex.Amount, ex.Note, pq.Array(ex.Tags)); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	ex.ID = rowID
	return c.JSON(http.StatusOK, ex)

}
