package model_test

import (
	"testing"
	"todo/src/model"
)

func Test_Titleが一文字の時に作成できる(t *testing.T) {
	title, err := model.NewTitle("a")
	if err != nil {
		t.Fatal(err)
	}
	if title.AsGoString() != "a" {
		t.Fatal("title is invalid")
	}
	t.Log("title is valid")
}

func Test_Titleが100文字の時に作成できない(t *testing.T) {
	word := "a"
	for i := 0; i < 99; i++ {
		word += "a"
	}
	_, err := model.NewTitle(word)
	if err != nil {
		t.Fatal("title is valid")
	}
	t.Log("title is invalid")
}

func Test_Titleが空文字の時に作成できない(t *testing.T) {
	_, err := model.NewTitle("")
	if err != nil {
		t.Fatal("title is valid")
	}
	t.Log("title is invalid")
}

func Test_CompleteをToggleすると値が反転する(t *testing.T) {
	completed_false := model.NewCompleted(false)
	completed_true := completed_false.Toggle()
	if completed_true.AsGoBool() != true {
		t.Fatal("completed is invalid")
	}
	t.Log("completed is valid")
}

// func Test_Todoを作成する(t *testing.T) {
// 	title, _ := model.NewTitle("title")
// 	completed := model.NewCompleted(false)
// 	lastUpdate := model.NewLastUpdate(model.Time{}.Now().AsGoTime())
// 	createdAt := model.NewCreatedAt(model.Time{}.Now().AsGoTime())
// 	todo := model.NewTodo(title, completed, lastUpdate, createdAt)
// 	if todo.Title.AsGoString() != "title" {
// 		t.Fatal("title is invalid")
// 	}
// 	if todo.Completed.AsGoBool() != false {
// 		t.Fatal("completed is invalid")
// 	}
// 	if todo.LastUpdate.AsGoString() != model.Time{}.Now().AsGoString() {
// 		t.Fatal("lastUpdate is invalid")
// 	}
// 	if todo.CreatedAt.AsGoString() != model.Time{}.Now().AsGoString() {
// 		t.Fatal("createdAt is invalid")
// 	}
// 	t.Log("todo is valid")
// }