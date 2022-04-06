package repositories

import (
	"database/sql"
	"errors"
	"shopingList/pkg/models"
)

type RefbookCategoriesRepository struct {
	db models.DB
}

func NewRefbookCategoriesRepository(db models.DB) RefbookCategoriesRepository {
	if db == nil {
		panic("db param is nil")
	}

	return RefbookCategoriesRepository{db: db}
}

func (s *RefbookCategoriesRepository) GetAll() ([]models.RefbookCategory, error) {
	rows, err := s.db.Query(
		`SELECT id,
       			title, 
       			UNIX_TIMESTAMP(created_at), 
				UNIX_TIMESTAMP(updated_at)
		FROM ` + models.RefbookCategoryTableName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound{}
		}

		return nil, err
	}
	defer rows.Close()

	var items []models.RefbookCategory

	for rows.Next() {
		var m models.RefbookCategory
		err = rows.Scan(&m.ID, &m.Title, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}

		items = append(items, m)
	}

	return items, nil
}

func (s *RefbookCategoriesRepository) Clear() error {
	_, err := s.db.Exec(`DELETE FROM ` + models.RefbookCategoryTableName)
	if err != nil {
		errors.New("Error delete all from category; " + err.Error())
	}

	return nil
}

func (s *RefbookCategoriesRepository) Create(form *models.RefbookCategory) (int64, error) {
	result, err := s.db.Exec(`INSERT INTO `+models.RefbookCategoryTableName+` (
                    title) 
		VALUES (?)`,
		form.Title)

	if err != nil {
		return 0, errors.New("Error insert category; " + err.Error())
	}

	return result.LastInsertId()
}
