package infra_test

import (
	"context"
	"testing"
	"time"
	"todo/internal/di"
	"todo/internal/domain"
	infra "todo/internal/infra"

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

func NewTestTodoRepository(t *testing.T) (infra.TodoRepositorier, error) {
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
	text , err := domain.NewTitle("test")
	if err != nil {
		t.Fatal(err)
	}
	todo := domain.NewTodo(
		0,
		*text,
		domain.NewCompleted(false),
		domain.NewLastUpdate(domain.NewDomainTime(time.Now())),
		domain.NewCreatedAt(domain.NewDomainTime(time.Now())),
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
	text1 , err := domain.NewTitle("test")
	if err != nil {
		t.Fatal(err)
	}
	todo1 := domain.NewTodo(
		0,
		*text1,
		domain.NewCompleted(false),
		domain.NewLastUpdate(domain.NewDomainTime(time.Now())),
		domain.NewCreatedAt(domain.NewDomainTime(time.Now())),
	)
	err = repository.Create(todo1)
	if err != nil {
		t.Fatal(err)
	}

	text2 , err := domain.NewTitle("test")
	if err != nil {
		t.Fatal(err)
	}
	todo2 := domain.NewTodo(
		0,
		*text2,
		domain.NewCompleted(false),
		domain.NewLastUpdate(domain.NewDomainTime(time.Now())),
		domain.NewCreatedAt(domain.NewDomainTime(time.Now())),
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
