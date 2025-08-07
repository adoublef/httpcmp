package httpecho

import (
	"cmp"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func Echo() *echo.Echo {
	e := echo.New()

	e.GET("/ping", handlePing())
	e.GET("/users/:user/books/:book", handleParameters())

	return e
}

func handlePing() echo.HandlerFunc {
	type response struct {
		Val string `json:"value"`
	}
	// no autocomplete
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, response{"pong"})
	}
}

func handleParameters() echo.HandlerFunc {
	type response struct {
		User int `json:"user"`
		Book int `json:"book"`
	}
	parse := func(c echo.Context) (user, book int, err error) {
		u, err1 := strconv.Atoi(c.Param("user"))
		b, err2 := strconv.Atoi(c.Param("book"))
		return u, b, cmp.Or(err1, err2)
	}
	return func(c echo.Context) error {
		user, book, err := parse(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, response{user, book})
	}
}
