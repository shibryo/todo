package repository_test

import (
	"context"
	"testing"
	"time"
	"todo/internal/di"
	"todo/internal/model"
	"todo/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupDB(t *testing.T) string{
	// Setup
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(45*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	return connStr
}

func NewTestTodoRepository(t *testing.T) (repository.TodoRepositorier, error) {
	dsn := SetupDB(t)
	repository, err := di.NewDITodoRepository(dsn)
	if err != nil {
		return nil, err
	}
	return repository, nil
}
func TestTodoRepository_FindAll(t *testing.T) {
	t.Parallel()
	repository, err := NewTestTodoRepository(t)
	if err != nil {
		t.Fatal(err)
	}

	// set data
	text , err := model.NewTitle("test")
	if err != nil {
		t.Fatal(err)
	}
	todo := model.NewTodo(
		0,
		text,
		model.NewCompleted(false),
		model.NewLastUpdate(model.NewModelTime(time.Now())),
		model.NewCreatedAt(model.NewModelTime(time.Now())),
	)
	err = repository.Create(todo)
	if err != nil {
		t.Fatal(err)
	}

	got, err := repository.FindAll()
	if err != nil {
		t.Fatal(err)
	}
	gotCount := len(got)
	assert.Equal(t, 1, gotCount)
}

func Test_Todoが2件取得できる(t *testing.T) {
	t.Parallel()
	repository, err := NewTestTodoRepository(t)
	if err != nil {
		t.Fatal(err)
	}

	// set 2data
	text1 , err := model.NewTitle("test")
	if err != nil {
		t.Fatal(err)
	}
	todo1 := model.NewTodo(
		0,
		text1,
		model.NewCompleted(false),
		model.NewLastUpdate(model.NewModelTime(time.Now())),
		model.NewCreatedAt(model.NewModelTime(time.Now())),
	)
	err = repository.Create(todo1)
	if err != nil {
		t.Fatal(err)
	}

	text2 , err := model.NewTitle("test")
	if err != nil {
		t.Fatal(err)
	}
	todo2 := model.NewTodo(
		0,
		text2,
		model.NewCompleted(false),
		model.NewLastUpdate(model.NewModelTime(time.Now())),
		model.NewCreatedAt(model.NewModelTime(time.Now())),
	)
	err = repository.Create(todo2)
	if err != nil {
		t.Fatal(err)
	}

	got, err := repository.FindAll()
	if err != nil {
		t.Fatal(err)
	}
	gotCount := len(got)
	assert.Equal(t, 2, gotCount)
}

func Test_Todoが更新できる(t *testing.T) {
	t.Parallel()
	repository, err := NewTestTodoRepository(t)
	if err != nil {
		t.Fatal(err)
	}

	// set data
	text , err := model.NewTitle("test")
	if err != nil {
		t.Fatal(err)
	}
	todo := model.NewTodo(
		0,
		text,
		model.NewCompleted(false),
		model.NewLastUpdate(model.NewModelTime(time.Now())),
		model.NewCreatedAt(model.NewModelTime(time.Now())),
	)
	err = repository.Create(todo)
	if err != nil {
		t.Fatal(err)
	}

	updateTodo, err := repository.FindByID(1)
	if err != nil {
		t.Fatal(err)
	}
	updateTitle, err := model.NewTitle("update")
	updateTodo.UpdateTitle(updateTitle)
	if err != nil {
		t.Fatal(err)
	}
	err = repository.Update(updateTodo)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "update", updateTodo.Title.AsGoString())
}


func Test_Todoが1件削除できる(t *testing.T) {
	t.Parallel()
	repository, err := NewTestTodoRepository(t)
	if err != nil {
		t.Fatal(err)
	}

	// set data
	text , err := model.NewTitle("test")
	if err != nil {
		t.Fatal(err)
	}
	todo := model.NewTodo(
		1,
		text,
		model.NewCompleted(false),
		model.NewLastUpdate(model.NewModelTime(time.Now())),
		model.NewCreatedAt(model.NewModelTime(time.Now())),
	)
	err = repository.Create(todo)
	if err != nil {
		t.Fatal(err)
	}

	deleteTodo, err := repository.FindByID(1)
	if err != nil {
		t.Fatal(err)
	}
	deletableID := &model.DeletableID{ID: deleteTodo.ID}
	err = repository.Delete(deletableID)
	if err != nil {
		t.Fatal(err)
	}
	got, err := repository.FindAll()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(got))
}
