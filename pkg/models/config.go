package models

import (
	"errors"
)

// AppConfig - base app config structure
type AppConfig struct {
	Server                  ServerConfig       `json:"server"`
	Database                DatabaseConfig     `json:"database"`
	SmsEnabled              bool               `json:"smsEnabled"`
	SmsAeroConfig           SmsAeroConfig      `json:"smsAeroConfig"`
	FirebaseCredentialsFile string             `json:"firebaseCredentialsFile"`
	RedisConfig             RedisConfig        `json:"redisConfig"`
	LoginLimiterConfig      LoginLimiterConfig `json:"loginLimiter"`
	LogLevel                string             `json:"logLevel"`
	TelegramBotToken        string             `json:"tgBotToken"`
	DebugPhones             []int64            `json:"debugPhones"`
}

// ServerConfig - server config
type ServerConfig struct {
	Port string `json:"port"`
}

// DatabaseConfig database config structure
type DatabaseConfig struct {
	Address  string `json:"address"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"dbName"`
}

type SmsAeroConfig struct {
	Email    string `json:"email"`
	ApiKey   string `json:"apiKey"`
	Sign     string `json:"sign"`
	Channel  string `json:"channel"`
	TestMode bool   `json:"testMode"`
}

func (s *SmsAeroConfig) Validate() (bool, error) {
	if s.Email == "" {
		return false, errors.New("email is empty")
	}

	if s.ApiKey == "" {
		return false, errors.New("apiKey is empty")
	}

	if s.Sign == "" {
		return false, errors.New("sign is empty")
	}

	if s.Channel == "" {
		return false, errors.New("channel is empty")
	}

	return true, nil
}

func (s *AppConfig) HasFirebaseCredentials() bool {
	return s.FirebaseCredentialsFile != ""
}

type RedisConfig struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type LoginLimiterConfig struct {
	SeqLimitSeconds int `json:"seqLimitSeconds"`
	DailyLimitCount int `json:"dailyLimitCount"`
}
