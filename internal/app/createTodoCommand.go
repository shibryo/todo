package app

import (
	"net/http"
	"todo/internal/domain"

	"github.com/labstack/echo/v4"
)

func (t *TodoComandService)CreateTodoCommand(todo *domain.Todo, c echo.Context) (bool, error) {
	err := t.repository.Create(todo)
	if err != nil {
		return true, c.JSON(http.StatusInternalServerError, err)
	}
	return false, nil
}