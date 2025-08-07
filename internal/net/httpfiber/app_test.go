// https://docs.gofiber.io/api/app#test
package httpfiber_test

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	. "github.com/adoublef/httpcmp/internal/net/httpfiber"
	"github.com/gofiber/fiber/v3"
	"go.adoublef.dev/testing/is"
)

func Test_handlePing(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		app := newApp(t)

		req, err := http.NewRequest("GET", "/ping", nil)
		is.OK(t, err)

		// Tets calls [httputil.DumpRequest] which buffers
		// request into memory. This is not ideal for testing when large payloads
		// are expected (file uploads).
		//
		// Config is optional and defaults to a 1s error on timeout,
		// 0 timeout will disable it completely.
		//
		// we are testing the route, not the server
		//
		// failing to set the TestConfig will use the default 1s timeout.
		res, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		is.OK(t, err) // GET /
		is.Equal(t, res.StatusCode, http.StatusOK)

		var data struct {
			Val string `json:"value"`
		}
		is.OK(t, json.NewDecoder(res.Body).Decode(&data))
		is.Equal(t, data.Val, "pong")

		is.OK(t, res.Body.Close())
	})
}

func Test_handleParameters(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		app := newApp(t)

		req, err := http.NewRequest("GET", "/users/1/books/1", nil)
		is.OK(t, err)

		// failing to set the TestConfig will use the default 1s timeout.
		res, err := app.Test(req, fiber.TestConfig{Timeout: 0})
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
		app := newApp(t)

		req, err := http.NewRequest("GET", "/users/a/books/1", nil)
		is.OK(t, err)

		res, err := app.Test(req, fiber.TestConfig{Timeout: 0})
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

func Test_handleUpload(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		app := newApp(t)

		// create a multipart
		pr, pw := io.Pipe()
		defer pr.Close()

		mw := multipart.NewWriter(pw)
		go func() {
			defer pw.Close()
			defer mw.Close()

			w, err := mw.CreateFormFile("file", "hello.txt")
			if err != nil {
				pw.CloseWithError(err)
				return
			}
			_, err = w.Write([]byte("hello, world!\n"))
			if err != nil {
				pw.CloseWithError(err)
				return
			}
		}()

		req, err := http.NewRequest("POST", "/upload", pr)
		is.OK(t, err)

		req.Header.Set("Content-Type", mw.FormDataContentType())

		res, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		is.OK(t, err)
		is.Equal(t, res.StatusCode, http.StatusBadRequest)

		var data struct {
			N int `json:"bytesWritten"`
		}
		is.OK(t, json.NewDecoder(res.Body).Decode(&data))
		is.Equal(t, data.N, 14)

		is.OK(t, res.Body.Close())
	})
}

func newApp(t testing.TB) *fiber.App {
	t.Helper()

	app := App()
	t.Cleanup(func() { app.Shutdown() })

	return app
}
