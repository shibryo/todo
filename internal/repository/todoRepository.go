package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"todo/internal/model"

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
	FindAll() ([]*model.Todo, error)
	FindByID(id uint64) (*model.Todo, error)
	Create(todoModel *model.Todo) error
	Update(todoModel *model.Todo) error
	Delete(todoModel *model.Todo) error
}

// TodoRepositoryはTodoのリポジトリ構造体です。
type TodoRepository struct {
	db *bun.DB
}

// FindAllは全てのTodoを取得します。
func (t *TodoRepository) FindAll() ([]*model.Todo, error) {
	var todos []*Todo
	err := t.db.NewSelect().Model(&todos).Scan(context.TODO())
	if err != nil {
		slog.Error("todos are not found", "err", err, "todos", todos)
		return nil, fmt.Errorf("todos are not found: %w", err)
	}
	todoModels, err := convertToTodoModels(todos)
	if err != nil {
		return nil, fmt.Errorf("todos are invalid: %w", err)
	}
	return todoModels, nil
}

func convertToTodoModels(todos []*Todo) ([]*model.Todo, error) {
	todoModels := make([]*model.Todo, 0, len(todos))
	for _, todo := range todos {
		id := model.NewID(todo.ID)
		title, err := model.NewTitle(todo.Title)
		if err != nil {
			slog.Error("title is invalid","err", err, "todo", todo)
			return nil, fmt.Errorf("title is invalid: %w", err)
		}
		completed := model.NewCompleted(todo.Completed)
		lastUpdate := model.NewLastUpdate(model.NewLastUpdate(model.NewModelTime(todo.LastUpdate)))
		createdAt := model.NewCreatedAt(model.NewCreatedAt(model.NewModelTime(todo.CreatedAt)))
		todoModel := model.NewTodo(id, title, completed, lastUpdate, createdAt)
		todoModels = append(todoModels, todoModel)
	}
	return todoModels, nil
}

// FindByIDはIDを指定してTodoを取得します。
func (t *TodoRepository) FindByID(id uint64) (*model.Todo, error) {
	todo := new(Todo)
	err := t.db.NewSelect().Model(todo).Where("id = ?", id).Scan(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("todo is not found: %w", err, "id", id)
	}
	todoModel, err := convertToTodoModel(todo)
	if err != nil {
		return nil, fmt.Errorf("todo is invalid: %w", err, "todo", todo)
	}
	return todoModel, nil
}

func convertToTodoModel(todo *Todo) (*model.Todo, error) {
	id := model.NewID(todo.ID)
	title, err := model.NewTitle(todo.Title)
	if err != nil {
		slog.Error("title is invalid","err", err, "todo", todo)
		return nil, fmt.Errorf("title is invalid: %w", err)
	}
	completed := model.NewCompleted(todo.Completed)
	lastUpdate := model.NewLastUpdate(model.NewLastUpdate(model.NewModelTime(todo.LastUpdate)))
	createdAt := model.NewCreatedAt(model.NewCreatedAt(model.NewModelTime(todo.CreatedAt)))
	todoModel := model.NewTodo(id, title, completed, lastUpdate, createdAt)
	return todoModel, nil
}

// CreateはTodoを作成します。
func (t *TodoRepository) Create(todoModel *model.Todo) error {
	todo := convertToTodo(todoModel)
	_, err := t.db.NewInsert().Model(todo).Exec(context.TODO())
	if err != nil {
		slog.Error("todo is invalid", "err",err, "todo", todo)
		return fmt.Errorf("todo is invalid: %w", err)
	}
	return nil
}

func convertToTodo(todoModel *model.Todo) *Todo {
	return &Todo{
		ID:         uint64(todoModel.ID),
		Title:      todoModel.Title.AsGoString(),
		Completed:  todoModel.Completed.AsGoBool(),
		LastUpdate: todoModel.LastUpdate.AsGoTime(),
		CreatedAt:  todoModel.CreatedAt.AsGoTime(),
	}
}

// UpdateはTodoを更新します。
func (t *TodoRepository) Update(todoModel *model.Todo) error {
	todo := convertToTodo(todoModel)
	_, err := t.db.NewUpdate().Model(todo).WherePK().Exec(context.TODO())
	if err != nil {
		slog.Error("todo is invalid", "err",err, "todo", todo)
		return fmt.Errorf("todo is invalid: %w", err)
	}
	return nil
}

// DeleteはTodoを削除します。
func (t *TodoRepository) Delete(todoModel *model.Todo) error {
	todo := convertToTodo(todoModel)
	_, err := t.db.NewDelete().Model(todo).WherePK().Exec(context.TODO())
	if err != nil {
		slog.Error("todo is invalid", "err",err, "todo", todo)
		return fmt.Errorf("todo is invalid: %w", err)
	}
	return nil
}

func NewTodoRepository(db *bun.DB) TodoRepositorier {
	return &TodoRepository{db: db}
}
