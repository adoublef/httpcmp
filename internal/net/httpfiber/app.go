package httpfiber

import (
	"cmp"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func App() *fiber.App {
	app := fiber.New(fiber.Config{
		// We set a function once. Custom logic needed to handle json errors.
		// which is no different to the stdlib.
		// https://docs.gofiber.io/guide/error-handling/#custom-error-handler
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return ctx.Status(code).JSON(map[string]string{"message": e.Message})
		},
	})
	// adding middleware is not type-safe
	// app.Use()
	app.Get("/ping", handlePing())
	app.Get("/users/:user/books/:book", handleParameters())
	return app
}

func handlePing() fiber.Handler {
	type response struct {
		Val string `json:"value"`
	}
	// no autocomplete
	return func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(response{"pong"})
	}
}

func handleParameters() fiber.Handler {
	type response struct {
		User int `json:"user"`
		Book int `json:"book"`
	}
	parse := func(c fiber.Ctx) (user, book int, err error) {
		u, err1 := strconv.Atoi(c.Params("user"))
		b, err2 := strconv.Atoi(c.Params("book"))
		return u, b, cmp.Or(err1, err2)
	}
	return func(c fiber.Ctx) error {
		user, book, err := parse(c)
		if err != nil {
			// according to this doc (Jun 17, 2025) we either:
			// - panic and use a custom recover handler, but loose type of error
			// - use a custom type, no different to stdlib approach.
			// https://leapcell.io/blog/effective-error-handling-in-go-gin-framework
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(response{user, book})
	}
}
