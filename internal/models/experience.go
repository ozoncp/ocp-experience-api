package models

import (
	"fmt"
	"time"
)

// Experience describes user experience info
type Experience struct {
	Id     uint64
	UserId uint64
	Type   uint64
	From   time.Time
	To     time.Time
	Level  uint64
}

// NewExperience creates Experience new instance
func NewExperience(id, userId, t uint64, from, to time.Time, level uint64) Experience {
	return Experience{
		Id:     id,
		UserId: userId,
		Type:   t,
		From:   from,
		To:     to,
		Level:  level,
	}
}

// String method converts Experience instance to string representation
func (e *Experience) String() string {
	return fmt.Sprintf("Experience{Id : %v, UserId: %v, Type : %v, From : %v, To : %v, Level : %v}",
		e.Id, e.UserId, e.Type, e.From.String(), e.To.String(), e.Level)
}
