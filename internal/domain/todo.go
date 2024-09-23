package domain

import (
	"fmt"
	"time"
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

func NewCompleted(value bool) Completed {
	return Completed{value: value}
}

func (c Completed) AsGoBool() bool {
	return c.value
}

func (c Completed) Toggle() Completed {
	return NewCompleted(!c.value)
}

func (c Completed) ToCompleted() Completed {
	return NewCompleted(true)
}

func (c Completed) ToIncompleted() Completed {
	return NewCompleted(false)
}

type LastUpdate struct {
	date DomainTimer
}

func NewLastUpdate(date DomainTimer) LastUpdate {
	return LastUpdate{date: date}
}

func (l LastUpdate) AsGoString() string {
	return l.date.AsGoString()
}

func (l LastUpdate) AsGoTime() time.Time {
	return l.date.AsGoTime()
}

type CreatedAt struct {
	date DomainTimer
}

func NewCreatedAt(date DomainTimer) CreatedAt {
	return CreatedAt{date: date}
}

func (c CreatedAt) AsGoString() string {
	return c.date.AsGoString()
}

func (c CreatedAt) AsGoTime() time.Time {
	return c.date.AsGoTime()
}

type Todo struct {
	ID         ID 	   `bun:"id"`
	Title      Title       `bun:"title"`
	Completed  Completed   `bun:"completed"`
	LastUpdate LastUpdate  `bun:"last_update"`
	CreatedAt  CreatedAt   `bun:"created_at"`
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
		LastUpdate: NewLastUpdate(NewDomainTimeNow()),
		CreatedAt:  NewCreatedAt(NewDomainTimeNow()),
	}
}

func (t *Todo) Delete() *DeletableTodo {
	return NewDeletableTodo(t.ID)
}

func (t *Todo) UpdateTitle(title *Title) {
	t.Title = *title
	t.UpdateLastUpdate(NewDomainTimeNow())
}

func (t *Todo) ToggleCompleted() {
	t.Completed = t.Completed.Toggle()
	t.UpdateLastUpdate(NewDomainTimeNow())
}

func (t *Todo) ToCompleted() {
	t.Completed = t.Completed.ToCompleted()
	t.UpdateLastUpdate(NewDomainTimeNow())
}

func (t *Todo) ToIncompleted() {
	t.Completed = t.Completed.ToIncompleted()
	t.UpdateLastUpdate(NewDomainTimeNow())
}

func (t *Todo) UpdateCompleted(completed Completed) {
	t.Completed = completed
	t.UpdateLastUpdate(NewDomainTimeNow())
}

func (t *Todo) UpdateLastUpdate(now DomainTimer) {
	t.LastUpdate = NewLastUpdate(now)
}