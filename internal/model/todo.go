package model

import (
	"fmt"
	"log/slog"
	"time"
)

type ID uint64

func NewID(id uint64) ID {
	return ID(id)
}

type Title struct {
	text string
}

func NewTitle(text string) (*Title, error) {
	if len(text) == 0 {
		slog.Error("title is empty", "text", text)
		return nil, fmt.Errorf("title is empty")
	}
	if len(text) > 100 {
		slog.Error("title is too long", "text", text, "len", len(text))
		return nil, fmt.Errorf("title is too long: %d", len(text))
	}

	return &Title{text: text}, nil
}

func (t Title) AsGoString() string {
	return t.text
}

type Completed struct {
	value bool
}

func NewCompleted(value bool) *Completed {
	return &Completed{value: value}
}

func (c Completed) AsGoBool() bool {
	return c.value
}

func (c Completed) Toggle() Completed {
	return *NewCompleted(!c.value)
}

func (c Completed) ToCompleted() *Completed {
	return NewCompleted(true)
}

func (c Completed) ToIncompleted() *Completed {
	return NewCompleted(false)
}

type LastUpdate struct {
	date ModelTimer
}

func NewLastUpdate(date ModelTimer) *LastUpdate {
	return &LastUpdate{date: date}
}

func (l LastUpdate) AsGoString() string {
	return l.date.AsGoString()
}

func (l LastUpdate) AsGoTime() time.Time {
	return l.date.AsGoTime()
}

type CreatedAt struct {
	date ModelTimer
}

func NewCreatedAt(date ModelTimer) *CreatedAt {
	return &CreatedAt{date: date}
}

func (c CreatedAt) AsGoString() string {
	return c.date.AsGoString()
}

func (c CreatedAt) AsGoTime() time.Time {
	return c.date.AsGoTime()
}

// Todoモデル
type Todo struct {
	ID         ID 	   `bun:"id"`
	Title      Title       `bun:"title"`
	Completed  Completed   `bun:"completed"`
	LastUpdate LastUpdate  `bun:"last_update"`
	CreatedAt  CreatedAt   `bun:"created_at"`
}

// 削除専用のモデル
type DeletableID struct {
	ID ID 
}

func NewTodo(id ID, title *Title, completed *Completed, lastUpdate *LastUpdate, createdAt *CreatedAt) *Todo {
	return &Todo{
		ID:         id,
		Title:      *title,
		Completed:  *completed,
		LastUpdate: *lastUpdate,
		CreatedAt:  *createdAt,
	}
}

func (t *Todo) UpdateTitle(title *Title) {
	t.Title = *title
}

func (t *Todo) ToggleCompleted() {
	t.Completed = t.Completed.Toggle()
}

func (t *Todo) ToCompleted() {
	t.Completed = *t.Completed.ToCompleted()
}

func (t *Todo) ToIncompleted() {
	t.Completed = *t.Completed.ToIncompleted()
}

func (t *Todo) UpdateLastUpdate(now ModelTimer) {
	t.LastUpdate = *NewLastUpdate(now)
}

func (t * Todo) Delete() *DeletableID {
	return &DeletableID{ID: t.ID}
}