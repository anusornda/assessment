package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github/anusornda/assessment/expense"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	expense.InitDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(checkUserAuth)

	e.POST("expenses", expense.CreateExpensesHandler)
	e.GET("/expenses/:id", expense.GetExpensesByIdHandler)
	e.PUT("/expenses/:id", expense.UpdateExpensesByIdHandler)
	e.GET("/expenses", expense.GetExpensesHandler)

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))

	go func() {
		if err := e.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}

func checkUserAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Values("Authorization")
		if auth[0] == "November 10, 2009" {
			return next(c)
		}
		return echo.ErrUnauthorized
	}
}
