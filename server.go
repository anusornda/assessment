package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github/anusornda/assessment/expense"
	"log"
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
