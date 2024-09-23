package di

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	app "todo/internal/app"
	infra "todo/internal/infra"
	view "todo/internal/view"
)

type config struct {
	DSN string `env:"DSN"`
}

func getEnv() (config, error) {
	err := godotenv.Load()
	if err != nil {
		return config{}, fmt.Errorf("failed to load .env: %w", err)
	}

	c := config{
		DSN: os.Getenv("DSN"),
	}

	slog.Info("getEnv", "c", c)

	return c, nil
}

// TodoRepositoryに依存性を注入します。
func NewDITodoRepository(dsn string) (*infra.TodoRepository, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	// create todos table
	_, err := db.NewCreateTable().
		Model((*infra.Todo)(nil)).
		IfNotExists().
		Exec(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	todoRepository := infra.NewTodoRepository(db)

	return todoRepository, nil
}

// NewDITodoControllerはTodoControllerを生成します。
func NewDITodoController() (*view.TodoController, error) {
	cnf, err := getEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to get env: %w", err)
	}

	dsn := cnf.DSN

	todoRepository, err := NewDITodoRepository(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo repository: %w", err)
	}

	todoCommandService, err := NewDITodoCommandService()
	if err != nil {
		return nil, fmt.Errorf("failed to get todo command service: %w", err)
	}

	c := view.NewTodoController(todoCommandService, todoRepository)

	return c, nil
}

// NewDITodoCommandServiceはTodoCommandServiceを生成します。
func NewDITodoCommandService() (*app.TodoComandServiceImpl, error) {
	cnf, err := getEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to get env: %w", err)
	}

	dsn := cnf.DSN

	todoRepository, err := NewDITodoRepository(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo repository: %w", err)
	}

	c := app.NewTodoCommandServiceImpl(todoRepository)

	return c, nil
}
