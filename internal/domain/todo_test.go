package domain_test

import (
	"testing"
	"time"
	"todo/internal/domain"

	"github.com/stretchr/testify/assert"
)

func Test_Titleが一文字の時に作成できる(t *testing.T) {
	t.Parallel()
	title, err := domain.NewTitle("a")

	assert.Nil(t, err)
	assert.Equal(t, title.AsGoString(), "a")
}

func Test_Titleが101文字の時に作成できない(t *testing.T) {
	t.Parallel()
	word := "a"
	for i := 0; i < 100; i++ {
		word += "a"
	}
	
	_, err := domain.NewTitle(word)

	assert.Equal(t, err.Error(), "title is too long")
}

func Test_Titleが空文字の時に作成できない(t *testing.T) {
	t.Parallel()
	_, err := domain.NewTitle("")

	assert.Equal(t, err.Error(), "title is empty")
}

func Test_CompleteをToggleすると値が反転する(t *testing.T) {
	t.Parallel()
	completed_false := domain.NewCompleted(false)
	completed_true := completed_false.Toggle()

	assert.Equal(t, completed_true.AsGoBool(), true)
}

func setupTodo() (*domain.Todo, domain.DomainTimer, ) {
	// ctrl := gomock.NewController(t)
	// mock_timer := model_mock.NewMockModelTimer(ctrl)

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

	assert.Equal(t, todo.ID, domain.NewID(1))
	assert.Equal(t, todo.Title.AsGoString(), "title")
	assert.Equal(t, todo.Completed.AsGoBool(), false)
	assert.Equal(t, todo.LastUpdate.AsGoTime(), now.AsGoTime())
	assert.Equal(t, todo.CreatedAt.AsGoTime(), now.AsGoTime())

}

func Test_TodoのTitleを更新できる(t *testing.T) {
	t.Parallel()
    todo, _ := setupTodo()
    newText := "new_title"
    newTitle, _ := domain.NewTitle(newText)

    todo.UpdateTitle(newTitle)

    if todo.Title.AsGoString() != newText {
        t.Fatal("title is invalid")
    }
    t.Log("title is valid")
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