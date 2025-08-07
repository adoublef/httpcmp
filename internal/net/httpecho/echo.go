package httpecho

import (
	"cmp"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func Echo() *echo.Echo {
	e := echo.New()
	e.GET("/ping", handlePing())
	e.GET("/users/:user/books/:book", handleParameters())
	e.POST("/upload", handleUpload(io.Discard))
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

func handleUpload(sink io.Writer) echo.HandlerFunc {
	// require the underlying [http.Request] for the multipart reader.
	return func(c echo.Context) error {
		mr, err := c.Request().MultipartReader()
		if err != nil {
			return echo.NewHTTPError(http.StatusUnsupportedMediaType, err.Error())
		}
		p, err := mr.NextPart()
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		defer p.Close()
		// todo: handle filename + formname
		_, err = io.Copy(sink, p)
		return err
	}
}
