package repositories

import (
	"errors"
	"shopingList/pkg/models"
)

type ListsRepository struct {
	DB models.DB
}

func (s *ListsRepository) UpdateList(list *models.List) error {
	_, err := s.DB.Exec(`UPDATE sl_item_list SET name=?, updated_at=FROM_UNIXTIME(?), received_at=FROM_UNIXTIME(?), 
				is_deleted=? 
				WHERE id=? AND owner_id=?`,
		list.Name, list.UpdatedAt, list.ReceivedAt, list.IsDeleted, list.ID, list.OwnerID)

	if err != nil {
		return errors.New("Error update list; " + err.Error())
	}

	return nil
}

func (s *ListsRepository) CreateList(list *models.List) error {
	_, err := s.DB.Exec(`INSERT INTO  sl_item_list (id, owner_id, name, is_template, created_at, updated_at, received_at, is_deleted) 
		VALUES (?, ?, ?, ?, FROM_UNIXTIME(?), FROM_UNIXTIME(?), FROM_UNIXTIME(?), ?)`,
		list.ID, list.OwnerID, list.Name, list.IsTemplate, list.CreatedAt, list.UpdatedAt, list.ReceivedAt, list.IsDeleted)

	if err != nil {
		return errors.New("Error insert list; " + err.Error())
	}

	return nil
}
