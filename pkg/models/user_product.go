package models

import (
	"github.com/asaskevich/govalidator"
)

type UserProduct struct {
	ID              string `json:"id" valid:"uuid,required" db:"id"`
	OwnerID         string `json:"owner_id" valid:"uuid,required" db:"owner_id"`
	CategoryID      int64  `json:"category_id" valid:"int,required" db:"category_id"`
	GlobalProductId int64  `json:"global_product_id" valid:"int" db:"global_product_id"`
	Name            string `json:"name" valid:"stringlength(1|100),required" db:"name"`
	CreatedAt       int64  `json:"created_at" valid:"int,required" db:"created_at"`
	UpdatedAt       int64  `json:"updated_at" valid:"int,required" db:"updated_at"`
	ReceivedAt      int64  `json:"received_at" db:"received_at"`
	IsDeleted       bool   `json:"is_deleted" valid:"required" db:"is_deleted"`
	IsFavorite      bool   `json:"is_favorite" valid:"required" db:"is_favorite"`
}

func (s *UserProduct) Validate() (bool, error) {
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *UserProduct) IsEqual(item *UserProduct) bool {
	return s.ID == item.ID &&
		s.OwnerID == item.OwnerID &&
		s.CategoryID == item.CategoryID &&
		s.GlobalProductId == item.GlobalProductId &&
		s.Name == item.Name &&
		s.UpdatedAt == item.UpdatedAt &&
		s.CreatedAt == item.CreatedAt &&
		s.IsDeleted == item.IsDeleted &&
		s.IsFavorite == item.IsFavorite
}
