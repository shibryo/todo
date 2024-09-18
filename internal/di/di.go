package di

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"todo/internal/controller"
	"todo/internal/repository"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type config struct {
	DSN string `env:"DSN"`
}

func getEnv()  (config, error) {
	err := godotenv.Load()
	if err != nil {
		slog.Error("failed to load .env file", "err", err)
		return config{}, fmt.Errorf("failed to load .env file: %w", err)
	}
	c := config{
		DSN: os.Getenv("DSN"),
	}
	slog.Debug("getEnv", "c", c)
	return c, nil
}

// TodoRepositoryに依存性を注入します。
func NewDITodoRepository(dsn string) (repository.TodoRepositorier, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	// create todos table
	_, err := db.NewCreateTable().
		Model((*repository.Todo)(nil)).
		IfNotExists().
		Exec(context.TODO())
	if err != nil {	
		slog.Error("failed to create table", "err", err)
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	todoRepository := repository.NewTodoRepository(db)
	return todoRepository, nil
}

// NewDITodoControllerはTodoControllerを生成します。
func NewDITodoController() (*controller.TodoController, error) {
	cnf, err := getEnv()
	if err != nil {
		slog.Error("failed to get env", "err", err)
		return nil, fmt.Errorf("failed to get env: %w", err)
	}
	dsn := cnf.DSN
	todoRepository, err := NewDITodoRepository(dsn)
	if err != nil {
		slog.Error("failed to create todo repository", "err", err)
		return nil, fmt.Errorf("failed to create todo repository: %w", err)
	}
	c := controller.NewTodoController(todoRepository)
	return c, nil
}