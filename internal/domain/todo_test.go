package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "todo/internal/domain"
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
	assert.Equal(t, "title is empty", err.Error())
}

func Test_CompleteをToggleすると値が反転する(t *testing.T) {
	t.Parallel()
	completedFalse := domain.NewCompleted(false)

	completedTrue := completedFalse.Toggle()
	assert.True(t, completedTrue.AsGoBool())
}

func setupTodo() (*domain.Todo, *domain.Time) {
	now := domain.NewDomainTime(time.Date(2024, 9, 4, 12, 0, 0, 0, time.Local))
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

func Test_TodoのCompleteをToggleできる(t *testing.T) {
	t.Parallel()
	todo, _ := setupTodo()

	todo.ToggleCompleted()

	if todo.Completed.AsGoBool() != true {
		t.Fatal("completed is invalid")
	}

	t.Log("completed is valid")
}

func Test_TodoのLastUpdateを更新できる(t *testing.T) {
	t.Parallel()
	todo, _ := setupTodo()
	now := domain.NewDomainTime(time.Date(2024, 9, 4, 12, 0, 0, 0, time.Local))

	todo.UpdateLastUpdate(now)

	if todo.LastUpdate.AsGoString() != "2024-09-04 12:00:00" {
		t.Fatal("lastUpdate is invalid", todo.LastUpdate.AsGoString())
	}
}
