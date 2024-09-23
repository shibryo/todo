package app_test

import (
	"context"
	"testing"
	"time"
	"todo/internal/app"
	"todo/internal/di"
	"todo/internal/domain"

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


func TestCreateTodoCommand_Todo作成が成功する(t *testing.T) {
	t.Parallel()
	dsn := SetupDB(t)
	repository, err := di.NewDITodoRepository(dsn)
	if err != nil {
		t.Fatal(err)
	}
	service := app.NewTodoCommandServiceImpl(repository)

	todoData := app.NewToDoData(1, "title", false)
	err = service.CreateTodoCommand(todoData)
	assert.NoError(t, err)
}



func TestUpdateTodoCommand_Todoを作成してからTodo更新が成功する(t *testing.T) {
	t.Parallel()
	dsn := SetupDB(t)
	repository, err := di.NewDITodoRepository(dsn)
	if err != nil {
		t.Fatal(err)
	}
	service := app.NewTodoCommandServiceImpl(repository)

	title, err := domain.NewTitle("old_title")
	if err != nil {
		t.Fatal(err)
	}
	id := uint64(1)
	oldTodo := domain.NewTodo(domain.NewID(id), *title, 
		domain.NewCompleted(false), 
		domain.NewLastUpdate(domain.NewDomainTime(time.Now())), 
		domain.NewCreatedAt(domain.NewDomainTime(time.Now())),
	)
	err = repository.Create(oldTodo)
	if err != nil {
		t.Fatal(err)
	}

	newTodoData := app.NewToDoData(1, "new_title", true)
	err = service.UpdateTodoCommand(newTodoData)

	assert.NoError(t, err)
	got, err := repository.FindByID(id)
	assert.NoError(t, err)
	assert.Equal(t, newTodoData.Title, got.Title.AsGoString())
}

func TestDeleteTodoCommand_Todoを作成してからTodo削除が成功する(t *testing.T) {
	t.Parallel()
	dsn := SetupDB(t)
	repository, err := di.NewDITodoRepository(dsn)
	if err != nil {
		t.Fatal(err)
	}
	service := app.NewTodoCommandServiceImpl(repository)

	title, err := domain.NewTitle("title")
	if err != nil {
		t.Fatal(err)
	}
	id := uint64(1)
	todo := domain.NewTodo(domain.NewID(id), *title,
		domain.NewCompleted(false),
		domain.NewLastUpdate(domain.NewDomainTime(time.Now())),
		domain.NewCreatedAt(domain.NewDomainTime(time.Now())),
	)
	err = repository.Create(todo)
	if err != nil {
		t.Fatal(err)
	}

	todoIDData := app.NewTodoIDData(id)
	err = service.DeleteTodoCommand(todoIDData)

	assert.NoError(t, err)
	_, err = repository.FindByID(id)
	assert.Error(t, err)
}

func TestFindAllCommand_Todoが2件取得できる(t *testing.T) {
	t.Parallel()
	dsn := SetupDB(t)
	repository, err := di.NewDITodoRepository(dsn)
	if err != nil {
		t.Fatal(err)
	}
	service := app.NewTodoCommandServiceImpl(repository)

	title1, err := domain.NewTitle("title1")
	if err != nil {
		t.Fatal(err)
	}
	todo1 := domain.NewTodo(domain.NewID(1), *title1,
		domain.NewCompleted(false),
		domain.NewLastUpdate(domain.NewDomainTime(time.Now())),
		domain.NewCreatedAt(domain.NewDomainTime(time.Now())),
	)
	err = repository.Create(todo1)
	if err != nil {
		t.Fatal(err)
	}
	title2, err := domain.NewTitle("title2")
	if err != nil {
		t.Fatal(err)
	}
	todo2 := domain.NewTodo(domain.NewID(2), *title2,
		domain.NewCompleted(false),
		domain.NewLastUpdate(domain.NewDomainTime(time.Now())),
		domain.NewCreatedAt(domain.NewDomainTime(time.Now())),
	)
	err = repository.Create(todo2)
	if err != nil {
		t.Fatal(err)
	}

	todos, err := service.FindAllCommand()
	assert.NoError(t, err)
	assert.Len(t, todos, 2)

}

func TestFindByIdCommand_IDを指定してTodoが取得できる(t *testing.T) {
	t.Parallel()
	dsn := SetupDB(t)
	repository, err := di.NewDITodoRepository(dsn)
	if err != nil {
		t.Fatal(err)
	}
	service := app.NewTodoCommandServiceImpl(repository)

	title, err := domain.NewTitle("title")
	if err != nil {
		t.Fatal(err)
	}
	id := uint64(1)
	todo := domain.NewTodo(domain.NewID(id), *title,
		domain.NewCompleted(false),
		domain.NewLastUpdate(domain.NewDomainTime(time.Now())),
		domain.NewCreatedAt(domain.NewDomainTime(time.Now())),
	)
	err = repository.Create(todo)
	if err != nil {
		t.Fatal(err)
	}

	todoIDData := app.NewTodoIDData(id)
	result, err := service.FindByIdCommand(todoIDData)
	assert.NoError(t, err)
	assert.Equal(t, id, result.ID.AsGoUint64())
	assert.Equal(t, title.AsGoString(), result.Title.AsGoString())
}