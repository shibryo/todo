package app

import "todo/internal/domain"

// TodoRepositorierはTodoのリポジトリインターフェースです。
type TodoRepositorier interface {
	FindAll() ([]*domain.Todo, error)
	FindByID(id uint64) (*domain.Todo, error)
	Create(todoModel *domain.Todo) error
	Update(todoModel *domain.Todo) error
	Delete(todoModel *domain.DeletableTodo) error
}