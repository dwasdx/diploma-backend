package models

import (
	"errors"
	"github.com/asaskevich/govalidator"
)

type List struct {
	ID         string `json:"id" valid:"uuid,required" db:"id"`
	OwnerID    string `json:"owner_id" valid:"uuid,required" db:"owner_id"`
	Name       string `json:"name" valid:"stringlength(1|100),required" db:"name"`
	IsTemplate bool   `json:"is_template" db:"is_template"`
	CreatedAt  int64  `json:"created_at" valid:"int,required" db:"created_at"`
	UpdatedAt  int64  `json:"updated_at" valid:"int,required" db:"updated_at"`
	ReceivedAt int64  `json:"received_at" db:"received_at"`
	IsDeleted  bool   `json:"is_deleted" valid:"required" db:"is_deleted"`
}

func (s *List) Validate() (bool, error) {
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

// IsEqual returns true if receiver is equal l2
func (l List) IsEqual(l2 List) bool {
	return l.ID == l2.ID &&
		l.OwnerID == l2.OwnerID &&
		l.Name == l2.Name &&
		l.IsTemplate == l2.IsTemplate &&
		l.CreatedAt == l2.CreatedAt &&
		l.UpdatedAt == l2.UpdatedAt &&
		l.IsDeleted == l2.IsDeleted
}
