package di

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	infra "todo/internal/infra"
	view "todo/internal/view"

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
		return config{}, err
	}
	c := config{
		DSN: os.Getenv("DSN"),
	}
	slog.Info("getEnv", "c", c)
	return c, nil
}

// TodoRepositoryに依存性を注入します。
func NewDITodoRepository(dsn string) (infra.TodoRepositorier, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	// create todos table
	_, err := db.NewCreateTable().
		Model((*infra.Todo)(nil)).
		IfNotExists().
		Exec(context.TODO())
	if err != nil {	
		return nil, err
	}
	todoRepository := infra.NewTodoRepository(db)
	return todoRepository, nil
}

// NewDITodoControllerはTodoControllerを生成します。
func NewDITodoController() (*view.TodoController, error) {
	cnf, err := getEnv()
	if err != nil {
		return nil, err
	}
	dsn := cnf.DSN
	todoRepository, err := NewDITodoRepository(dsn)
	if err != nil {
		return nil, err
	}
	c := view.NewTodoController(todoRepository)
	return c, nil
}