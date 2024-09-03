package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type TodoController struct {
}

func NewTodoController() *TodoController {
	c := &TodoController{}
	return c
}

// GetHello godoc
// @Summary Get hello
// @Description get hello
// @ID get-hello
// @Produce  plain
// @Success 200 {string} string "Hello, World!"
// @Router /v1/api/ [get]
func(ctrl *TodoController) GetHello() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}
	
}