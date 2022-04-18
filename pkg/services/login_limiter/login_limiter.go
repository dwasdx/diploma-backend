package login_limiter

import (
	"github.com/pkg/errors"
	"shopingList/api/controllers/users"
	"strconv"
)

type LoginLimiter struct {
	storage           LoginLimiterStorage
	seqLimitSeconds   int
	dailyLimitSeconds int
	dailyLimitCount   int
}

type LoginLimiterStorage interface {
	SaveSeqLimit(phone string, limitSeconds int) (int64, error)
	SaveDailyLimit(phone string, limitSeconds int) (DailyLimit, error)
	GetSeqLimit(phone string) (int64, error)
	GetDailyLimit(phone string) (DailyLimit, error)
}

type DailyLimit struct {
	Exist    bool
	NextTime int64
	Counter  int
}

func NewLoginLimiter(storage LoginLimiterStorage, seqLimitSeconds int, dailyLimitCount int) *LoginLimiter {
	return &LoginLimiter{
		storage:           storage,
		seqLimitSeconds:   seqLimitSeconds,
		dailyLimitSeconds: 86400,
		dailyLimitCount:   dailyLimitCount}
}

func (s *LoginLimiter) IsAllowLogin(phone int64) (bool, error) {
	phoneStr := strconv.FormatInt(phone, 10)

	// Проверка на частотный лимит
	nextTime, err := s.storage.GetSeqLimit(phoneStr)
	if err != nil {
		return false, errors.Wrap(err, "error is IsAllowLogin")
	}

	if nextTime != 0 {
		return false, nil
	}

	// Проверка на дневной лимит
	dailyLimit, err := s.storage.GetDailyLimit(phoneStr)
	if err != nil {
		return false, errors.Wrap(err, "error is IsAllowLogin")
	}

	if s.dailyLimitExceeded(dailyLimit) {
		return false, nil
	}

	return true, nil
}

func (s *LoginLimiter) RememberOperation(phone int64) (int64, error) {
	phoneStr := strconv.FormatInt(phone, 10)

	nextTime, err := s.storage.SaveSeqLimit(phoneStr, s.seqLimitSeconds)
	if err != nil {
		return 0, errors.Wrap(err, "error save sequence limit")
	}

	limit, err := s.storage.SaveDailyLimit(phoneStr, s.dailyLimitSeconds)
	if err != nil {
		return 0, errors.Wrap(err, "error save daily limit")
	}

	if s.dailyLimitExceeded(limit) {
		return limit.NextTime, nil
	}

	return nextTime, nil
}

func (s *LoginLimiter) NextTimeRequest(phone int64) (users.LoginLimiterNextTimeResponse, error) {
	phoneStr := strconv.FormatInt(phone, 10)
	response := users.LoginLimiterNextTimeResponse{}

	limit, err := s.storage.GetDailyLimit(phoneStr)
	if err != nil {
		return response, errors.Wrap(err, "Error getDailyLimit() in NextTimeRequest")
	}

	if s.dailyLimitExceeded(limit) {
		response.TypeLimit = users.NewTypeLoginLimit(users.TypeDailyLimit)
		response.SetNextTimeRequest(limit.NextTime)
		return response, nil
	}

	nextTime, err := s.storage.GetSeqLimit(phoneStr)
	if err != nil {
		return response, errors.Wrap(err, "Error getSeqLimit() in NextTimeRequest")
	}

	if nextTime != 0 {
		response.SetNextTimeRequest(nextTime)
		response.TypeLimit = users.NewTypeLoginLimit(users.TypeSequenceLimit)

		return response, nil
	}

	response.IsAllow = true

	return response, nil
}

func (s *LoginLimiter) dailyLimitExceeded(limit DailyLimit) bool {
	if limit.Exist && limit.Counter >= s.dailyLimitCount {
		return true
	}

	return false
}
