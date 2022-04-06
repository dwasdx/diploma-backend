package repositories

import (
	"database/sql"
	"errors"
	"shopingList/pkg/models"
)

type UserProductsRepository struct {
	db models.DB
}

func NewUserProductsRepository(db models.DB) UserProductsRepository {
	if db == nil {
		panic("db param is nil")
	}

	return UserProductsRepository{db: db}
}

func (s *UserProductsRepository) GetOneById(id string) (models.UserProduct, error) {
	db := s.db
	row := db.QueryRow(
		`SELECT id,
       			name,
       			owner_id,
       			category_id,
       			global_product_id,
       			UNIX_TIMESTAMP(created_at), 
       			UNIX_TIMESTAMP(updated_at), 
       			UNIX_TIMESTAMP(received_at),
       			is_favorite,
       			is_deleted 
		FROM sl_user_products
		WHERE id = ?`, id)

	var i models.UserProduct
	err := row.Scan(
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
		if err == sql.ErrNoRows {
			return models.UserProduct{}, ErrNotFound{}
		}

		return models.UserProduct{}, err
	}
	return i, nil
}

func (s *UserProductsRepository) Create(userProduct *models.UserProduct) error {
	db := s.db

	_, err := db.Exec(`INSERT INTO sl_user_products (
                    id, name, owner_id, category_id, global_product_id, is_deleted, is_favorite, created_at, updated_at, received_at
                    ) 
		VALUES (?, ?, ?, ?, ?, ?, ?, FROM_UNIXTIME(?), FROM_UNIXTIME(?), FROM_UNIXTIME(?))`,
		userProduct.ID, userProduct.Name, userProduct.OwnerID, userProduct.CategoryID, userProduct.GlobalProductId, userProduct.IsDeleted, userProduct.IsFavorite,
		userProduct.CreatedAt, userProduct.UpdatedAt, userProduct.ReceivedAt)

	if err != nil {
		return errors.New("Error insert list; " + err.Error())
	}

	return nil
}

func (s *UserProductsRepository) Update(userProduct *models.UserProduct) error {
	_, err := s.db.Exec(`UPDATE sl_user_products
		SET name=?, category_id=?, is_deleted=?, is_favorite=?, updated_at=FROM_UNIXTIME(?), received_at=FROM_UNIXTIME(?)
		WHERE id=?`,
		userProduct.Name, userProduct.CategoryID, userProduct.IsDeleted, userProduct.IsFavorite, userProduct.UpdatedAt, userProduct.ReceivedAt, userProduct.ID)

	if err != nil {
		return errors.New("Error update userProduct; " + err.Error())
	}

	return nil
}

func (s *UserProductsRepository) Delete(Id string) error {
	_, err := s.db.Exec(`DELETE FROM sl_user_products WHERE id=?`, Id)
	if err != nil {
		return err
	}

	return nil
}
