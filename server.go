package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github/anusornda/assessment/expense"
	"log"
	"os"
)

func main() {

	expense.InitDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//e.Use(middleware.BasicAuth(func(userName, password string, c echo.Context) (bool, error) {
	//
	//	fmt.Printf("Authorization : %v", c.Request().Header.Values("Authorization"))
	//	log.Printf("userName: %s | password: %s", userName, password)
	//	if userName == "expenseApi" && password == "123456" {
	//		return true, nil
	//	}
	//	return false, nil
	//}))
	e.Use(checkUserAuth)

	e.POST("expenses", expense.CreateExpensesHandler)
	e.GET("/expenses/:id", expense.GetExpensesByIdHandler)
	e.PUT("/expenses/:id", expense.UpdateExpensesByIdHandler)
	e.GET("/expenses", expense.GetExpensesHandler)

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))

	log.Printf("Server started at : 2565")
	log.Fatal(e.Start(os.Getenv("PORT")))
	log.Printf("bye bye!")
}

func checkUserAuth(next echo.HandlerFunc) echo.HandlerFunc {
	log.Printf("## checkUserAuth ##")
	return func(c echo.Context) error {

		log.Printf("Authorization : %v", c.Request().Header.Values("Authorization"))

		auth := c.Request().Header.Values("Authorization")
		if auth[0] == "November 10, 2009" {
			log.Printf("return next(c)")
			return next(c)
		}
		log.Printf("echo.ErrUnauthorized")
		// Even if middleware reaches here, it still execute the next route, why?
		return echo.ErrUnauthorized
	}
}
