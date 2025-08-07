package httpgin

import (
	"cmp"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Engine() *gin.Engine {
	e := gin.Default()
	e.GET("/ping", handleOK())
	e.GET("/users/:user/books/:book", handleParameters())
	return e
}

func handleOK() gin.HandlerFunc {
	type response struct {
		Val string `json:"value"`
	}
	// the default variable name iven by the IDE is `ctx`
	// which can be misleading to the lanuage wide accepted
	// type of [context.Context]
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, response{"pong"})
	}
}

func handleParameters() gin.HandlerFunc {
	type response struct {
		User int `json:"user"`
		Book int `json:"book"`
	}
	parse := func(ctx *gin.Context) (user, book int, err error) {
		u, err1 := strconv.Atoi(ctx.Params.ByName("user"))
		b, err2 := strconv.Atoi(ctx.Params.ByName("book"))
		return u, b, cmp.Or(err1, err2)
	}
	return func(ctx *gin.Context) {
		user, book, err := parse(ctx)
		if err != nil {
			// according to this doc (Jun 17, 2025) we either:
			// - panic and use a custom recover handler, but loose type of error
			// - use a custom type, no different to stdlib approach.
			// https://leapcell.io/blog/effective-error-handling-in-go-gin-framework
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, response{user, book})
	}
}
