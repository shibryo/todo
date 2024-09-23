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
	FindAllCommand() ([]*domain.Todo, error)
	FindByIdCommand(id uint64) (*domain.Todo, error)
}

// TodoComandServiceはTodoのコマンドサービス構造体です。
type TodoComandServiceImpl struct {
	repository repository.TodoRepositorier
}

// NewTodoCommandServiceImplはTodoのコマンドサービスを生成します。
func NewTodoCommandServiceImpl(repository repository.TodoRepositorier) TodoComandService {
	return &TodoComandServiceImpl{repository: repository}
}

// CreateTodoCommandはTodoを作成します。
func (t *TodoComandServiceImpl)CreateTodoCommand(todo *domain.Todo) error {
	err := t.repository.Create(todo)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	return nil
}

// UpdateTodoCommandはTodoを更新します。
func (t *TodoComandServiceImpl)UpdateTodoCommand(newTodo *domain.Todo) error {
	oldTodo, err := t.repository.FindByID(newTodo.ID.AsGoUint64())
	if err != nil {
		return fmt.Errorf("failed to find todo: %w", err)
	}

	oldTodo.UpdateTitle(&newTodo.Title)
	oldTodo.UpdateCompleted(&newTodo.Completed)

	err = t.repository.Update(oldTodo)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}
	return nil
}

// DeleteTodoCommandはTodoを削除します。
func (t *TodoComandServiceImpl)DeleteTodoCommand(todo *domain.Todo) error {
	err := t.repository.Delete(todo)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	return nil
}

// FindAllCommandは全てのTodoを取得します。
func (t *TodoComandServiceImpl)FindAllCommand() ([]*domain.Todo, error) {
	todos, err := t.repository.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to find all todos: %w", err)
	}
	return todos, nil
}

// FindByIdCommandはIDを指定してTodoを取得します。
func (t *TodoComandServiceImpl)FindByIdCommand(id uint64) (*domain.Todo, error) {
	todo, err := t.repository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find todo by id: %w", err)
	}
	return todo, nil
}