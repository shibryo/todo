package repository

import (
	"context"
	"fmt"
	"time"
	"todo/internal/domain"

	"github.com/uptrace/bun"
)

// TodoはORMのモデルです。
type Todo struct {
	bun.BaseModel `bun:"table:todos"`
	ID         uint64 `bun:"id,pk,autoincrement,unique,notnull"`
	Title      string `bun:"title,notnull"`
	Completed  bool  `bun:"completed,notnull,default:false"`
	LastUpdate time.Time `bun:"last_update,notnull,default:current_timestamp"`
	CreatedAt  time.Time `bun:"created_at,notnull,default:current_timestamp"`
}

// TodoRepositorierはTodoのリポジトリインターフェースです。
type TodoRepositorier interface {
	FindAll() ([]*domain.Todo, error)
	FindByID(id uint64) (*domain.Todo, error)
	Create(todoModel *domain.Todo) error
	Update(todoModel *domain.Todo) error
	Delete(todoModel *domain.Todo) error
}

// TodoRepositoryはTodoのリポジトリ構造体です。
type TodoRepository struct {
	db *bun.DB
}

// FindAllは全てのTodoを取得します。
func (t *TodoRepository) FindAll() ([]*domain.Todo, error) {
	var todos []*Todo
	err := t.db.NewSelect().Model(&todos).Scan(context.TODO())
	if err != nil {
		return nil, err
	}
	todoModels, err := convertToTodoModels(todos)
	if err != nil {
		return nil, err
	}
	return todoModels, nil
}

func convertToTodoModels(todos []*Todo) ([]*domain.Todo, error) {
	todoModels := make([]*domain.Todo, 0, len(todos))
	for _, todo := range todos {
		id := domain.NewID(todo.ID)
		title, err := domain.NewTitle(todo.Title)
		if err != nil {
			return nil, fmt.Errorf("title is invalid: %w", err)
		}
		completed := domain.NewCompleted(todo.Completed)
		lastUpdate := domain.NewLastUpdate(domain.NewLastUpdate(domain.NewModelTime(todo.LastUpdate)))
		createdAt := domain.NewCreatedAt(domain.NewCreatedAt(domain.NewModelTime(todo.CreatedAt)))
		todoModel := domain.NewTodo(id, title, completed, lastUpdate, createdAt)
		todoModels = append(todoModels, todoModel)
	}
	return todoModels, nil
}

// FindByIDはIDを指定してTodoを取得します。
func (t *TodoRepository) FindByID(id uint64) (*domain.Todo, error) {
	todo := new(Todo)
	err := t.db.NewSelect().Model(todo).Where("id = ?", id).Scan(context.TODO())
	if err != nil {
		return nil, err
	}
	todoModel, err := convertToTodoModel(todo)
	if err != nil {
		return nil, err
	}
	return todoModel, nil
}

func convertToTodoModel(todo *Todo) (*domain.Todo, error) {
	id := domain.NewID(todo.ID)
	title, err := domain.NewTitle(todo.Title)
	if err != nil {
		return nil, fmt.Errorf("title is invalid: %w", err)
	}
	completed := domain.NewCompleted(todo.Completed)
	lastUpdate := domain.NewLastUpdate(domain.NewLastUpdate(domain.NewModelTime(todo.LastUpdate)))
	createdAt := domain.NewCreatedAt(domain.NewCreatedAt(domain.NewModelTime(todo.CreatedAt)))
	todoModel := domain.NewTodo(id, title, completed, lastUpdate, createdAt)
	return todoModel, nil
}

// CreateはTodoを作成します。
func (t *TodoRepository) Create(todoModel *domain.Todo) error {
	todo := convertToTodo(todoModel)
	_, err := t.db.NewInsert().Model(todo).Exec(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

func convertToTodo(todoModel *domain.Todo) *Todo {
	return &Todo{
		ID:         uint64(todoModel.ID),
		Title:      todoModel.Title.AsGoString(),
		Completed:  todoModel.Completed.AsGoBool(),
		LastUpdate: todoModel.LastUpdate.AsGoTime(),
		CreatedAt:  todoModel.CreatedAt.AsGoTime(),
	}
}

// UpdateはTodoを更新します。
func (t *TodoRepository) Update(todoModel *domain.Todo) error {
	todo := convertToTodo(todoModel)
	_, err := t.db.NewUpdate().Model(todo).WherePK().Exec(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

// DeleteはTodoを削除します。
func (t *TodoRepository) Delete(todoModel *domain.Todo) error {
	todo := convertToTodo(todoModel)
	_, err := t.db.NewDelete().Model(todo).WherePK().Exec(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

func NewTodoRepository(db *bun.DB) TodoRepositorier {
	return &TodoRepository{db: db}
}
