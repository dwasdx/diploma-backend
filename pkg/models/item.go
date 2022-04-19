package models

import (
	"errors"
	"github.com/asaskevich/govalidator"
)

type ListItem struct {
	ID         string     `json:"id" valid:"uuid,required"`
	Name       string     `json:"name" valid:"stringlength(1|140),required"`
	Value      string     `json:"value" valid:"stringlength(0|50),required"`
	IsMarked   bool       `json:"is_marked" valid:"required"`
	UserMarked NullString `json:"user_marked"`
	ListID     string     `json:"list_id" valid:"uuid,required"`
	CreatedAt  int64      `json:"created_at" valid:"int,required"`
	UpdatedAt  int64      `json:"updated_at" valid:"int,required"`
	ReceivedAt int64      `json:"received_at"`
	IsDeleted  bool       `json:"is_deleted" valid:"required"`
}

func (s *ListItem) Validate() (bool, error) {
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		return false, err
	}

	if s.IsMarked && s.UserMarked.IsEmpty() {
		return false, errors.New("UserMarked is empty")
	}

	if !s.UserMarked.IsEmpty() && !govalidator.IsUUID(s.UserMarked.String) {
		return false, errors.New("UserMarked format is wrong")
	}

	if s.CreatedAt == 0 {
		return false, errors.New("CreatedAt is 0")
	}

	if s.UpdatedAt == 0 {
		return false, errors.New("UpdatedAt is 0")
	}

	return true, nil
}

func (s *ListItem) IsEqualItem(item *ListItem) bool {
	return s.ID == item.ID &&
		s.ListID == item.ListID
}

func (s *ListItem) IsEqual(item *ListItem) bool {
	return s.ID == item.ID &&
		s.ListID == item.ListID &&
		s.Name == item.Name &&
		s.Value == item.Value &&
		s.IsMarked == item.IsMarked &&
		s.UserMarked.String == item.UserMarked.String &&
		s.UpdatedAt == item.UpdatedAt &&
		s.IsDeleted == item.IsDeleted
}
