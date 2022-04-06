package repositories

import (
	"database/sql"
	"errors"
	"shopingList/pkg/models"
)

type ItemsRepository struct {
	db models.DB
}

func NewItemsRepository(db models.DB) ItemsRepository {
	if db == nil {
		panic("db param is nil")
	}

	return ItemsRepository{db: db}
}

// Item return shopping list item with specified id
func (s *ItemsRepository) GetItem(id string, listId string) (models.ListItem, error) {
	db := s.db
	row := db.QueryRow(
		`SELECT id,
       			name, 
       			value,
       			is_marked, 
       			user_marked_id,
       			list_id, 
       			is_deleted, 
       			UNIX_TIMESTAMP(created_at), 
       			UNIX_TIMESTAMP(updated_at)
		FROM sl_item
		WHERE id = ? AND list_id =?`, id, listId)

	var i models.ListItem
	err := row.Scan(
		&i.ID, &i.Name, &i.Value, &i.IsMarked, &i.UserMarked, &i.ListID, &i.IsDeleted, &i.CreatedAt, &i.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.ListItem{}, ErrNotFound{}
		}

		return models.ListItem{}, err
	}
	return i, nil
}

func (s *ItemsRepository) CreateItem(item *models.ListItem) error {
	v := item.UserMarked.SqlValue()

	_, err := s.db.Exec(`INSERT INTO sl_item (
                    id, name, value, is_marked, user_marked_id, list_id, is_deleted, created_at, updated_at, received_at
                    ) 
		VALUES (?, ?, ?, ?, ?, ?, ?, FROM_UNIXTIME(?), FROM_UNIXTIME(?), FROM_UNIXTIME(?))`,
		item.ID, item.Name, item.Value, item.IsMarked, v, item.ListID, item.IsDeleted,
		item.CreatedAt, item.UpdatedAt, item.ReceivedAt)

	if err != nil {
		return errors.New("Error insert item; " + err.Error())
	}

	return nil
}

func (s *ItemsRepository) UpdateItem(item *models.ListItem) error {
	_, err := s.db.Exec(`UPDATE sl_item 
		SET name=?, value=?, is_marked=?, user_marked_id=?, is_deleted=?, updated_at=FROM_UNIXTIME(?), received_at=FROM_UNIXTIME(?)
		WHERE id=? AND list_id=?`,
		item.Name, item.Value, item.IsMarked, item.UserMarked.SqlValue(), item.IsDeleted, item.UpdatedAt, item.ReceivedAt,
		item.ID, item.ListID)

	if err != nil {
		return errors.New("Error update item; " + err.Error())
	}

	return nil
}
