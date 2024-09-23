package app

import (
	"fmt"
	"todo/internal/domain"
	repository "todo/internal/infra"
)

// TodoComandServiceはTodoのコマンドサービスインターフェースです。
type TodoComandService interface {
	CreateTodoCommand(todo *domain.Todo) error
	UpdateTodoCommand(todo *domain.Todo) error
	DeleteTodoCommand(todo *domain.Todo) error
}

// TodoComandServiceはTodoのコマンドサービス構造体です。
type TodoComandServiceImpl struct {
	repository repository.TodoRepositorier
}


func NewTodoCommandServiceImpl(repository repository.TodoRepositorier) TodoComandService {
	return &TodoComandServiceImpl{repository: repository}
}


func (t *TodoComandServiceImpl)CreateTodoCommand(todo *domain.Todo) error {
	err := t.repository.Create(todo)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	return nil
}


func (t *TodoComandServiceImpl)UpdateTodoCommand(todo *domain.Todo) error {
	err := t.repository.Update(todo)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}
	return nil
}

func (t *TodoComandServiceImpl)DeleteTodoCommand(todo *domain.Todo) error {
	err := t.repository.Delete(todo)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	return nil
}