package app

import (
	"net/http"
	"todo/internal/domain"
	repository "todo/internal/infra"

	"github.com/labstack/echo/v4"
)

// TodoComandServiceはTodoのコマンドサービスインターフェースです。
type TodoComandService interface {
	CreateTodoCommand(todo *domain.Todo, c echo.Context) (bool, error)
}

// TodoComandServiceはTodoのコマンドサービス構造体です。
type TodoComandServiceImpl struct {
	repository repository.TodoRepositorier
}


func NewTodoCommandServiceImpl(repository repository.TodoRepositorier) TodoComandService {
	return &TodoComandServiceImpl{repository: repository}
}


func (t *TodoComandServiceImpl)CreateTodoCommand(todo *domain.Todo, c echo.Context) (bool, error) {
	err := t.repository.Create(todo)
	if err != nil {
		return true, c.JSON(http.StatusInternalServerError, err)
	}
	return false, nil
}

