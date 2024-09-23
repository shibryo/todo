package app

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=app_mock

import (
	"fmt"

	domain "todo/internal/domain"
)

type TodoIDData struct {
	ID uint64 `json:"id"`
}

type TodoData struct {
	ID        uint64 `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func NewToDoData(id uint64, title string, completed bool) TodoData {
	return TodoData{
		ID:        id,
		Title:     title,
		Completed: completed,
	}
}

func NewTodoIDData(id uint64) TodoIDData {
	return TodoIDData{
		ID: id,
	}
}

// TodoComandServiceはTodoのコマンドサービスインターフェースです。
type TodoComandService interface {
	CreateTodoCommand(todo TodoData) error
	UpdateTodoCommand(todo TodoData) error
	DeleteTodoCommand(todo TodoIDData) error
	FindAllCommand() ([]*domain.Todo, error)
	FindByIDCommand(todo TodoIDData) (*domain.Todo, error)
}

// TodoComandServiceはTodoのコマンドサービス構造体です。
type TodoComandServiceImpl struct {
	repository TodoRepositorier
}

// NewTodoCommandServiceImplはTodoのコマンドサービスを生成します。
func NewTodoCommandServiceImpl(repository TodoRepositorier) *TodoComandServiceImpl {
	return &TodoComandServiceImpl{repository: repository}
}

// CreateTodoCommandはTodoを作成します。
func (t *TodoComandServiceImpl) CreateTodoCommand(todoData TodoData) error {
	title, err := domain.NewTitle(todoData.Title)
	if err != nil {
		return fmt.Errorf("title is invalid: %w", err)
	}

	err = t.repository.Create(domain.Create(
		domain.NewID(todoData.ID),
		*title,
		domain.NewCompleted(todoData.Completed),
	))
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	return nil
}

// UpdateTodoCommandはTodoを更新します。
func (t *TodoComandServiceImpl) UpdateTodoCommand(newTodo TodoData) error {
	oldTodo, err := t.repository.FindByID(newTodo.ID)
	if err != nil {
		return fmt.Errorf("failed to find todo: %w", err)
	}

	newTitle, err := domain.NewTitle(newTodo.Title)
	if err != nil {
		return fmt.Errorf("title is invalid: %w", err)
	}

	oldTodo.UpdateTitle(newTitle)
	oldTodo.UpdateCompleted(domain.NewCompleted(newTodo.Completed))

	err = t.repository.Update(oldTodo)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	return nil
}

// DeleteTodoCommandはTodoを削除します。
func (t *TodoComandServiceImpl) DeleteTodoCommand(id TodoIDData) error {
	oldTodo, err := t.repository.FindByID(id.ID)
	if err != nil {
		return fmt.Errorf("failed to find todo: %w", err)
	}

	deletableTodo := domain.NewDeletableTodo(oldTodo.ID)

	err = t.repository.Delete(deletableTodo)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	return nil
}

// FindAllCommandは全てのTodoを取得します。
func (t *TodoComandServiceImpl) FindAllCommand() ([]*domain.Todo, error) {
	todos, err := t.repository.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to find all todos: %w", err)
	}

	return todos, nil
}

// FindByIDCommandはIDを指定してTodoを取得します。
func (t *TodoComandServiceImpl) FindByIDCommand(reqTodo TodoIDData) (*domain.Todo, error) {
	todo, err := t.repository.FindByID(reqTodo.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find todo by id: %w", err)
	}

	return todo, nil
}
