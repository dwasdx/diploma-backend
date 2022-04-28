package users

import (
	"log"
)

const TypeDailyLimit = "daily"
const TypeSequenceLimit = "sequence"

type TypeLimit string

func NewTypeLoginLimit(typeLimit string) TypeLimit {
	if typeLimit != TypeDailyLimit && typeLimit != TypeSequenceLimit {
		log.Fatal("Error wrong type limit")
	}

	return TypeLimit(typeLimit)
}

type LoginLimiter interface {
	IsAllowLogin(phone int64) (bool, error)
	RememberOperation(phone int64) (nextTime int64, error error)
	NextTimeRequest(phone int64) (LoginLimiterNextTimeResponse, error)
}

type LoginLimiterNextTimeResponse struct {
	IsAllow         bool      `json:"is_allow"`
	NextTimeRequest int64     `json:"next_time_request"`
	TypeLimit       TypeLimit `json:"type_limit"`
}

func (s *LoginLimiterNextTimeResponse) SetNextTimeRequest(nextTime int64) {
	if nextTime != 0 {
		s.IsAllow = false
		s.NextTimeRequest = nextTime
	} else {
		s.IsAllow = true
		s.NextTimeRequest = 0
	}
}

func (s *LoginLimiterNextTimeResponse) SetTimeLimit(typeLimit TypeLimit) {
	s.IsAllow = false
	s.TypeLimit = typeLimit
}
