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

	h := expense.InitDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(expense.CheckUserAuth())

	e.POST("expenses", h.CreateExpensesHandler)
	e.GET("/expenses/:id", h.GetExpensesByIdHandler)
	e.PUT("/expenses/:id", h.UpdateExpensesByIdHandler)
	e.GET("/expenses", h.GetExpensesHandler)

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
