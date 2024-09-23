package domain

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=app_mock

import (
	"errors"
	"time"
)

var (
	ErrTitleTooShort = errors.New("title is too short")
	ErrTitleTooLong  = errors.New("title is too long")
)

type ID uint64

func NewID(id uint64) ID {
	return ID(id)
}

func (i ID) AsGoUint64() uint64 {
	return uint64(i)
}

type Title struct {
	text string
}

func NewTitle(text string) (*Title, error) {
	minTextLength := 1
	if len(text) < minTextLength {
		return nil, ErrTitleTooShort
	}

	maxTextLength := 100
	if len(text) > maxTextLength {
		return nil, ErrTitleTooLong
	}

	return &Title{text: text}, nil
}

func (t Title) AsGoString() string {
	return t.text
}

type Completed struct {
	value bool
}

func NewCompleted(value bool) Completed {
	return Completed{value: value}
}

func (c Completed) AsGoBool() bool {
	return c.value
}

type LastUpdate struct {
	date Timer
}

func NewLastUpdate(date Timer) LastUpdate {
	return LastUpdate{date: date}
}

func (l LastUpdate) AsGoString() string {
	return l.date.AsGoString()
}

func (l LastUpdate) AsGoTime() time.Time {
	return l.date.AsGoTime()
}

type CreatedAt struct {
	date Timer
}

func NewCreatedAt(date Timer) CreatedAt {
	return CreatedAt{date: date}
}

func (c CreatedAt) AsGoString() string {
	return c.date.AsGoString()
}

func (c CreatedAt) AsGoTime() time.Time {
	return c.date.AsGoTime()
}

type Todo struct {
	ID         ID         `bun:"id"`
	Title      Title      `bun:"title"`
	Completed  Completed  `bun:"completed"`
	LastUpdate LastUpdate `bun:"last_update"`
	CreatedAt  CreatedAt  `bun:"created_at"`
}

func NewTodo(id ID, title Title, completed Completed, lastUpdate LastUpdate, createdAt CreatedAt) *Todo {
	return &Todo{
		ID:         id,
		Title:      title,
		Completed:  completed,
		LastUpdate: lastUpdate,
		CreatedAt:  createdAt,
	}
}

type DeletableTodo struct {
	ID ID `bun:"id"`
}

func NewDeletableTodo(id ID) *DeletableTodo {
	return &DeletableTodo{ID: id}
}

func Create(id ID, title Title, completed Completed) *Todo {
	return &Todo{
		ID:         id,
		Title:      title,
		Completed:  completed,
		LastUpdate: NewLastUpdate(TimerNow()),
		CreatedAt:  NewCreatedAt(TimerNow()),
	}
}

func (t *Todo) Delete() *DeletableTodo {
	return NewDeletableTodo(t.ID)
}

func (t *Todo) UpdateTitle(title *Title) {
	t.Title = *title
	t.UpdateLastUpdate(TimerNow())
}

func (t *Todo) UpdateCompleted(completed Completed) {
	t.Completed = completed
	t.UpdateLastUpdate(TimerNow())
}

func (t *Todo) UpdateLastUpdate(now Timer) {
	t.LastUpdate = NewLastUpdate(now)
}
