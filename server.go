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

	e.Use(middleware.BasicAuth(func(userName, password string, c echo.Context) (bool, error) {

		log.Printf("userName: %s | password: %s", userName, password)
		if userName == "expenseApi" && password == "123456" {
			return true, nil
		}
		return false, nil
	}))

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))

	log.Printf("Server started at : 2565")
	log.Fatal(e.Start(os.Getenv("PORT")))
	log.Printf("bye bye!")
}
