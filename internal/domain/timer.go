package domain

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=domain_mock

import "time"

type Timer interface {
	AsGoString() string
	AsGoTime() time.Time
}

type Time struct {
	date time.Time
}

func NewTime(date time.Time) *Time {
	return &Time{date: date}
}

func (t Time) AsGoString() string {
	return t.date.Format("2006-01-02 15:04:05")
}

func (t Time) AsGoTime() time.Time {
	return t.date
}

func TimerNow() *Time {
	return NewTime(time.Now())
}
