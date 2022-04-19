package models

import (
	"errors"
	"github.com/asaskevich/govalidator"
)

// ListShare store model
type ListShare struct {
	ID         string `json:"id" valid:"uuid,required"`
	ListID     string `json:"list_id" valid:"uuid,required"`
	ToUserID   string `json:"to_user_id" valid:"uuid,required"`
	OwnerID    string `json:"owner_id" valid:"uuid,required"`
	Status     int    `json:"status" valid:"range(0|2),required"`
	CreatedAt  int64  `json:"created_at" valid:"int,required"`
	UpdatedAt  int64  `json:"updated_at" valid:"int,required"`
	ReceivedAt int64  `json:"received_at"`
	IsDeleted  bool   `json:"is_deleted" valid:"required"`
}

const (
	ShareStatusNew      int = 0
	ShareStatusAccepted int = 1
	ShareStatusRefused  int = 2
)

func (s *ListShare) Validate() (bool, error) {
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		return false, err
	}

	if s.CreatedAt == 0 {
		return false, errors.New("CreatedAt is 0")
	}

	if s.UpdatedAt == 0 {
		return false, errors.New("UpdatedAt is 0")
	}

	return true, nil
}

// IsEqual returns true if receiver is equal s2
func (s ListShare) IsEqual(s2 ListShare) bool {
	return s.ID == s2.ID &&
		s.ListID == s2.ListID &&
		s.ToUserID == s2.ToUserID &&
		s.Status == s2.Status &&
		s.CreatedAt == s2.CreatedAt &&
		s.UpdatedAt == s2.UpdatedAt &&
		s.IsDeleted == s2.IsDeleted
}
