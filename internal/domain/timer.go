package domain

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=model_mock

import "time"

type DomainTimer interface {
	AsGoString() string
	AsGoTime() time.Time
}

type DomainTime struct {
	date time.Time
}

func NewDomainTime(date time.Time) DomainTimer {
	return &DomainTime{date: date}
}

func (t DomainTime) AsGoString() string {
	return t.date.Format("2006-01-02 15:04:05")
}

func (t DomainTime) AsGoTime() time.Time {
	return t.date
}

func NewDomainTimeNow() DomainTimer {
	return NewDomainTime(time.Now())
}