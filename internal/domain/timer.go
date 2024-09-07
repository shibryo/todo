package domain

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=model_mock

import "time"

type ModelTimer interface {
	AsGoString() string
	AsGoTime() time.Time
}

type ModelTime struct {
	date time.Time
}

func NewModelTime(date time.Time) ModelTimer {
	return &ModelTime{date: date}
}

func (t ModelTime) AsGoString() string {
	return t.date.Format("2006-01-02 15:04:05")
}

func (t ModelTime) AsGoTime() time.Time {
	return t.date
}

func NewModelTimeNow() ModelTimer {
	return NewModelTime(time.Now())
}