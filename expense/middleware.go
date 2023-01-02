package expense

import "github.com/labstack/echo/v4"

func CheckUserAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Values("Authorization")
			if auth != nil && auth[0] == "November 10, 2009" {
				return next(c)
			}
			return echo.ErrUnauthorized
		}
	}
}
