package model

import (
	"fmt"
	"time"
)

type ID uint64
type Title struct {
	text string
}
func NewTitle(text string) (*Title, error) {
	if len(text) == 0 {
		return nil, fmt.Errorf("title is empty")
	}
	if len(text) > 100 {
		return nil, fmt.Errorf("title is too long")
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
type LastUpdate struct {
	date time.Time
}
func NewLastUpdate(date time.Time) *LastUpdate {
	return &LastUpdate{date: date}
}
func (l LastUpdate) AsGoString() string {
	return l.date.Format("2006-01-02 15:04:05")
}
func (l LastUpdate) AsGoTime() time.Time {
	return l.date
}
type CreatedAt struct {
	date time.Time
}
func NewCreatedAt(date time.Time) *CreatedAt {
	return &CreatedAt{date: date}
}

func (t *Todo) UpdateLastUpdate(lastUpdate *LastUpdate) {
	t.LastUpdate = *NewLastUpdate(time.Now())
}
func (c CreatedAt) AsGoString() string {
	return c.date.Format("2006-01-02 15:04:05")
}
func (c CreatedAt) AsGoTime() time.Time {
	return c.date
}

type Todo struct {
	ID         ID 	   `bun:"id"`
	Title      Title       `bun:"title"`
	Completed  Completed   `bun:"completed"`
	LastUpdate LastUpdate  `bun:"last_update"`
	CreatedAt  CreatedAt   `bun:"created_at"`
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
