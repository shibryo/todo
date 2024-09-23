package domain_test

import (
	"testing"
	"time"

	domain "todo/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Titleが一文字の時に作成できる(t *testing.T) {
	t.Parallel()
	title, err := domain.NewTitle("a")
	require.NoError(t, err)

	assert.Equal(t, "a", title.AsGoString())
}

func Test_Titleが101文字の時に作成できない(t *testing.T) {
	t.Parallel()

	word := "a"
	for range [101]int{} {
		word += "a"
	}

	_, err := domain.NewTitle(word)

	assert.Equal(t, "title is too long", err.Error())
}

func Test_Titleが空文字の時に作成できない(t *testing.T) {
	t.Parallel()
	_, err := domain.NewTitle("")
	assert.Equal(t, "title is too short", err.Error())
}

func setupTodo() (*domain.Todo, *domain.Time) {
	now := domain.NewTime(time.Date(2024, 9, 4, 12, 0, 0, 0, time.Local))
	text := "title"
	id := domain.NewID(1)
	title, _ := domain.NewTitle(text)
	completed := domain.NewCompleted(false)
	lastUpdate := domain.NewLastUpdate(now)
	createdAt := domain.NewCreatedAt(now)
	todo := domain.NewTodo(id, *title, completed, lastUpdate, createdAt)

	return todo, now
}

func Test_Todo(t *testing.T) {
	t.Parallel()

	todo, now := setupTodo()

	assert.Equal(t, domain.NewID(1), todo.ID)
	assert.Equal(t, "title", todo.Title.AsGoString())
	assert.False(t, todo.Completed.AsGoBool())
	assert.Equal(t, now.AsGoTime(), todo.LastUpdate.AsGoTime())
	assert.Equal(t, now.AsGoTime(), todo.CreatedAt.AsGoTime())
}

func Test_TodoのTitleを更新できる(t *testing.T) {
	t.Parallel()
	todo, _ := setupTodo()

	newText := "new_title"

	newTitle, err := domain.NewTitle(newText)
	require.NoError(t, err)

	todo.UpdateTitle(newTitle)
	assert.Equal(t, "new_title", todo.Title.AsGoString())
}

func Test_TodoのCompletedを更新できる(t *testing.T) {
	t.Parallel()
	todo, _ := setupTodo()

	todo.UpdateCompleted(domain.NewCompleted(true))

	assert.True(t, todo.Completed.AsGoBool())
}

func Test_TodoのLastUpdateを更新できる(t *testing.T) {
	t.Parallel()
	todo, _ := setupTodo()
	now := domain.NewTime(time.Date(2024, 9, 4, 12, 0, 0, 0, time.Local))

	todo.UpdateLastUpdate(now)

	if todo.LastUpdate.AsGoString() != "2024-09-04 12:00:00" {
		t.Fatal("lastUpdate is invalid", todo.LastUpdate.AsGoString())
	}
}

func Test_TodoのDeleteでDeletableTodoが作成される(t *testing.T) {
	t.Parallel()
	todo, _ := setupTodo()

	deletableTodo := todo.Delete()

	if deletableTodo.ID.AsGoUint64() != 1 {
		t.Fatal("deletableTodo is invalid", deletableTodo.ID.AsGoUint64())
	}
}

func Test_TodoのCreateでTodoが作成される(t *testing.T) {
	t.Parallel()
	id := domain.NewID(1)
	title, _ := domain.NewTitle("title")
	completed := domain.NewCompleted(false)

	todo := domain.Create(id, *title, completed)

	assert.Equal(t, domain.NewID(1), todo.ID)
	assert.Equal(t, "title", todo.Title.AsGoString())
	assert.False(t, todo.Completed.AsGoBool())
	assert.NotEqual(t, "", todo.LastUpdate.AsGoString())
	assert.NotEqual(t, "", todo.CreatedAt.AsGoString())
}
