package repositories

import (
	"database/sql"
	"errors"
	"shopingList/pkg/models"
)

type RefbookProductsRepository struct {
	db models.DB
}

func NewRefbookProductsRepository(db models.DB) RefbookProductsRepository {
	if db == nil {
		panic("db param is nil")
	}

	return RefbookProductsRepository{db: db}
}

func (s *RefbookProductsRepository) GetAll() (*[]models.RefbookProduct, error) {
	rows, err := s.db.Query(
		`SELECT id,
       			title, 
				category_id,
       			UNIX_TIMESTAMP(created_at), 
				UNIX_TIMESTAMP(updated_at)
		FROM ` + models.RefbookProductTableName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound{}
		}

		return nil, err
	}
	defer rows.Close()

	var list []models.RefbookProduct

	for rows.Next() {
		var m models.RefbookProduct
		err = rows.Scan(&m.ID, &m.Title, &m.CategoryId, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}

		list = append(list, m)
	}

	return &list, nil
}

func (s *RefbookProductsRepository) Clear() error {
	_, err := s.db.Exec(`DELETE FROM ` + models.RefbookProductTableName)
	if err != nil {
		errors.New("Error delete all from products; " + err.Error())
	}

	return nil
}

func (s *RefbookProductsRepository) Create(form *models.RefbookProduct) (int64, error) {
	result, err := s.db.Exec(`INSERT INTO `+models.RefbookProductTableName+` (
                    title, category_id) 
		VALUES (?, ?)`,
		form.Title, form.CategoryId)

	if err != nil {
		return 0, errors.New("Error insert product; " + err.Error())
	}

	return result.LastInsertId()
}
