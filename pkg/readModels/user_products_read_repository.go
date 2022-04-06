package readModels

import (
	"database/sql"
	"shopingList/pkg/models"
	"shopingList/pkg/repositories"
)

type UserProductsReadRepository struct {
	useProducts map[string]*models.UserProduct
	db          *sql.DB
}

func NewUserProductsReadRepository(db *sql.DB) UserProductsReadRepository {
	if db == nil {
		panic("DB is nil")
	}

	return UserProductsReadRepository{db: db}
}

func (s *UserProductsReadRepository) GetUpdatedUserProductsForUser(ownerID string, receivedAt int64) ([]models.UserProduct, error) {
	var items []models.UserProduct
	db := s.db
	rows, err := db.Query(
		`SELECT id, 
       			name,
       			owner_id,
       			category_id,
       			global_product_id,
       			UNIX_TIMESTAMP(i.created_at), 
       			UNIX_TIMESTAMP(i.updated_at), 
       			UNIX_TIMESTAMP(i.received_at),
       			is_favorite,
       			is_deleted 
			FROM sl_user_products AS i
			WHERE owner_id =? AND i.received_at >= FROM_UNIXTIME(?)`,
		ownerID, receivedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repositories.ErrNotFound{}
		}

		return nil, err
	}
	defer rows.Close()

	items, err = s.scanItemRows(rows)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *UserProductsReadRepository) scanItemRows(rows *sql.Rows) ([]models.UserProduct, error) {
	var items []models.UserProduct

	for rows.Next() {
		var i models.UserProduct
		err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.OwnerID,
			&i.CategoryID,
			&i.GlobalProductId,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ReceivedAt,
			&i.IsFavorite,
			&i.IsDeleted,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}
