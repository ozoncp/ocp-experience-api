package models

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/ozoncp/ocp-experience-api/pkg/ocp-experience-api"
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

// ConvertExperienceToAPI converts model.Experience to desc.Experience
func ConvertExperienceToAPI(experience *Experience) *desc.Experience {
	return &desc.Experience{
		Id:     experience.Id,
		UserId: experience.UserId,
		Type:   experience.Type,
		From:   timestamppb.New(experience.From),
		To:     timestamppb.New(experience.To),
		Level:  experience.Level,
	}
}

// ConvertAPIToExperience converts desc.Experience to model.Experience
func ConvertAPIToExperience(experience *desc.Experience) Experience {
	return Experience{
		Id:     experience.Id,
		UserId: experience.UserId,
		Type:   experience.Type,
		From:   experience.From.AsTime(),
		To:     experience.To.AsTime(),
		Level:  experience.Level,
	}
}
