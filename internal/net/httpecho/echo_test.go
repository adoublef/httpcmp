// https://echo.labstack.com/docs/testing
package httpecho_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/adoublef/httpcmp/internal/net/httpecho"
	"go.adoublef.dev/testing/is"
)

func Test_handlePing(t *testing.T) {
	c, baseURL := newServer(t)
	// make a query
	res, err := c.Get(baseURL + "/ping")
	is.OK(t, err)
	is.Equal(t, res.StatusCode, http.StatusOK)

	var data struct {
		Val string `json:"value"`
	}
	is.OK(t, json.NewDecoder(res.Body).Decode(&data))
	is.Equal(t, data.Val, "pong")

	is.OK(t, res.Body.Close())
}

func Test_handleParameters(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		c, baseURL := newServer(t)

		res, err := c.Get(baseURL + "/users/1/books/1")
		is.OK(t, err)
		is.Equal(t, res.StatusCode, http.StatusOK)

		var data struct {
			User int `json:"user"`
			Book int `json:"book"`
		}
		is.OK(t, json.NewDecoder(res.Body).Decode(&data))
		is.Equal(t, data.User, 1)
		is.Equal(t, data.User, data.Book)

		is.OK(t, res.Body.Close())
	})

	t.Run("ErrBadRequest", func(t *testing.T) {
		c, baseURL := newServer(t)

		res, err := c.Get(baseURL + "/users/a/books/1")
		is.OK(t, err)
		is.Equal(t, res.StatusCode, http.StatusBadRequest)

		var data struct {
			Error string `json:"message"`
		}
		is.OK(t, json.NewDecoder(res.Body).Decode(&data))
		is.True(t, data.Error != "")

		is.OK(t, res.Body.Close())
	})
}

func newServer(t testing.TB) (*http.Client, string) {
	t.Helper()

	// works with the standard libary [httptest] package
	s := httptest.NewServer(Echo())
	t.Cleanup(func() { s.Close() })

	return s.Client(), s.URL
}
