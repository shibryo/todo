package di

import (
	"context"
	"database/sql"
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
	DNS string `env:"DNS"`
}

func getEnv()  (config, error) {
	err := godotenv.Load()
	if err != nil {
		return config{}, err
	}
	c := config{
		DNS: os.Getenv("DSN"),
	}
	slog.Info("getEnv", "c", c)
	return c, nil
}

// TodoRepositoryに依存性を注入します。
func NewDITodoRepository() (repository.TodoRepositorier, error) {
	c, err := getEnv()
	if err != nil {
		return nil, err
	}
	dsn := c.DNS
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	// create todos table
	_, err = db.NewCreateTable().
		Model((*repository.Todo)(nil)).
		IfNotExists().
		Exec(context.TODO())
	if err != nil {	
		return nil, err
	}
	todoRepository := repository.NewTodoRepository(db)
	return todoRepository, nil
}

// NewDITodoControllerはTodoControllerを生成します。
func NewDITodoController() (*controller.TodoController, error) {
	todoRepository, err := NewDITodoRepository()
	if err != nil {
		return nil, err
	}
	c := controller.NewTodoController(todoRepository)
	return c, nil
}