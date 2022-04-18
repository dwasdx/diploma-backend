package services

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type FixedCodeGenerator struct {
}

type AuthCodeGenerator struct {
}

func (s FixedCodeGenerator) GetCode() string {
	return "11111"
}

func (s AuthCodeGenerator) GetCode() string {
	start := 0
	end := 99999

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	n := start + r1.Intn(end-start+1)
	log.Info("code: ", n)

	return fmt.Sprintf("%05d", n)
}
