package app_test

import (
	"context"
	"testing"
	"time"

	app "todo/internal/app"
	di "todo/internal/di"
	domain "todo/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupDB(t *testing.T) string {
	t.Helper()
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
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	return connStr
}

func TestCreateTodoCommand_Todo作成が成功する(t *testing.T) {
	t.Parallel()
	dsn := SetupDB(t)

	repository, err := di.NewDITodoRepository(dsn)
	require.NoError(t, err)

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
		domain.NewLastUpdate(domain.NewTime(time.Now())),
		domain.NewCreatedAt(domain.NewTime(time.Now())),
	)

	err = repository.Create(oldTodo)
	require.NoError(t, err)

	newTodoData := app.NewToDoData(1, "new_title", true)
	err = service.UpdateTodoCommand(newTodoData)
	require.NoError(t, err)

	got, err := repository.FindByID(id)
	require.NoError(t, err)

	assert.Equal(t, newTodoData.Title, got.Title.AsGoString())
}

func TestDeleteTodoCommand_Todoを作成してからTodo削除が成功する(t *testing.T) {
	t.Parallel()
	dsn := SetupDB(t)
	repository, err := di.NewDITodoRepository(dsn)
	require.NoError(t, err)

	service := app.NewTodoCommandServiceImpl(repository)

	title, err := domain.NewTitle("title")
	require.NoError(t, err)

	id := uint64(1)
	todo := domain.NewTodo(domain.NewID(id), *title,
		domain.NewCompleted(false),
		domain.NewLastUpdate(domain.NewTime(time.Now())),
		domain.NewCreatedAt(domain.NewTime(time.Now())),
	)

	err = repository.Create(todo)
	require.NoError(t, err)

	todoIDData := app.NewTodoIDData(id)

	err = service.DeleteTodoCommand(todoIDData)
	require.NoError(t, err)

	_, err = repository.FindByID(id)
	require.Error(t, err)

	assert.Contains(t, err.Error(), "failed to get todo")
}

func TestFindAllCommand_Todoが2件取得できる(t *testing.T) {
	t.Parallel()
	dsn := SetupDB(t)

	repository, err := di.NewDITodoRepository(dsn)
	require.NoError(t, err)

	service := app.NewTodoCommandServiceImpl(repository)

	title1, err := domain.NewTitle("title1")
	require.NoError(t, err)

	todo1 := domain.NewTodo(domain.NewID(1), *title1,
		domain.NewCompleted(false),
		domain.NewLastUpdate(domain.NewTime(time.Now())),
		domain.NewCreatedAt(domain.NewTime(time.Now())),
	)

	err = repository.Create(todo1)
	require.NoError(t, err)

	title2, err := domain.NewTitle("title2")
	require.NoError(t, err)

	todo2 := domain.NewTodo(domain.NewID(2), *title2,
		domain.NewCompleted(false),
		domain.NewLastUpdate(domain.NewTime(time.Now())),
		domain.NewCreatedAt(domain.NewTime(time.Now())),
	)
	err = repository.Create(todo2)
	require.NoError(t, err)

	todos, err := service.FindAllCommand()
	require.NoError(t, err)
	assert.Len(t, todos, 2)
}

func TestFindByIdCommand_IDを指定してTodoが取得できる(t *testing.T) {
	t.Parallel()
	dsn := SetupDB(t)

	repository, err := di.NewDITodoRepository(dsn)
	require.NoError(t, err)

	service := app.NewTodoCommandServiceImpl(repository)

	title, err := domain.NewTitle("title")
	if err != nil {
		t.Fatal(err)
	}

	id := uint64(1)
	todo := domain.NewTodo(domain.NewID(id), *title,
		domain.NewCompleted(false),
		domain.NewLastUpdate(domain.NewTime(time.Now())),
		domain.NewCreatedAt(domain.NewTime(time.Now())),
	)

	err = repository.Create(todo)
	require.NoError(t, err)

	todoIDData := app.NewTodoIDData(id)
	result, err := service.FindByIDCommand(todoIDData)
	require.NoError(t, err)
	assert.Equal(t, id, result.ID.AsGoUint64())
	assert.Equal(t, title.AsGoString(), result.Title.AsGoString())
}
