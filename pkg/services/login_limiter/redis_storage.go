package login_limiter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const PrefixDailyLimit = "daily_lm_"
const PrefixSeqLimit = "seq_lm_"

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(address string, password string, db int) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	return &RedisStorage{client: client}
}

func secondsToDuration(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}

func getSeqKey(phone string) string {
	return PrefixSeqLimit + phone
}

func getDailyKey(phone string) string {
	return PrefixDailyLimit + phone
}

func (s *RedisStorage) GetSeqLimit(phone string) (int64, error) {
	rawValue, err := s.client.Get(context.Background(), getSeqKey(phone)).Result()

	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, errors.Wrap(err, "error get seqLimit")
	}

	nextTime, err := strconv.ParseInt(rawValue, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "error parse seqLimit")
	}

	return nextTime, nil
}

func (s *RedisStorage) SaveSeqLimit(phone string, limitSeconds int) (int64, error) {
	nextTime := time.Now().UTC().Add(secondsToDuration(limitSeconds)).Unix()

	err := s.client.Set(context.Background(), getSeqKey(phone), nextTime, secondsToDuration(limitSeconds)).Err()
	if err != nil {
		return 0, errors.Wrap(err, "error set seqLimit")
	}

	return nextTime, nil
}

func (s *RedisStorage) GetDailyLimit(phone string) (DailyLimit, error) {
	limit := DailyLimit{Exist: false}

	value, err := s.client.Get(context.Background(), getDailyKey(phone)).Result()

	if err == redis.Nil {
		return limit, nil
	} else if err != nil {
		return limit, errors.Wrap(err, "error get seqLimit")
	}

	err = json.Unmarshal([]byte(value), &limit)
	if err != nil {
		log.Fatal("Error parse dailyLimit from redis")
	}

	limit.Exist = true

	return limit, nil
}

func (s *RedisStorage) SaveDailyLimit(phone string, limitSeconds int) (DailyLimit, error) {
	limit, err := s.GetDailyLimit(phone)
	if err != nil {
		return limit, errors.Wrap(err, "Err get daily limit")
	}

	if !limit.Exist {
		limit.NextTime = time.Now().UTC().Unix()
		limit.Counter = 1
	} else {
		limit.Counter = limit.Counter + 1
	}

	dataJson, err := json.Marshal(limit)
	if err != nil {
		return limit, errors.Wrap(err, fmt.Sprintf("Error marshalling dailyLimit %v", limit))
	}

	err = s.client.Set(context.Background(), getDailyKey(phone), dataJson, secondsToDuration(limitSeconds)).Err()
	if err != nil {
		return limit, errors.Wrap(err, "error set daily limit")
	}

	return limit, nil
}
