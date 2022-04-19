package models

import "errors"

const RefbookCategoryTableName = "sl_categories"
const RefbookProductTableName = "sl_products"

type RefbookCategory struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (s *RefbookCategory) Validate() (bool, error) {
	if s.Title == "" {
		return false, errors.New("Title is empty")
	}

	return true, nil
}

type RefbookProduct struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	CategoryId int64  `json:"category_id"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

func (s *RefbookProduct) Validate() (bool, error) {
	if s.Title == "" {
		return false, errors.New("Title is empty")
	}

	if s.CategoryId == 0 {
		return false, errors.New("CategoryId is empty")
	}

	return true, nil
}
