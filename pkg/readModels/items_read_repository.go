package readModels

import (
	"database/sql"
	"shopingList/pkg/models"
	"shopingList/pkg/repositories"
)

type ItemsReadRepository struct {
	db *sql.DB
}

func NewItemsReadRepository(db *sql.DB) ItemsReadRepository {
	if db == nil {
		panic("DB is nil")
	}

	return ItemsReadRepository{db: db}
}

// Вернуть обновленные товары для списков пользователя
func (s *ItemsReadRepository) GetUpdatedItemsForUser(listOwnerID string, receivedAt int64) ([]models.ListItem, error) {
	var items []models.ListItem
	db := s.db
	rows, err := db.Query(
		`SELECT i.id, 
       			i.name, 
       			value,
       			is_marked, 
       			user_marked_id,
       			list_id, 
       			UNIX_TIMESTAMP(i.created_at), 
       			UNIX_TIMESTAMP(i.updated_at), 
       			UNIX_TIMESTAMP(i.received_at),
       			 i.is_deleted 
			FROM sl_item AS i
			LEFT JOIN sl_item_list AS l ON (i.list_id = l.id) 
			WHERE l.owner_id =? AND i.received_at >= FROM_UNIXTIME(?)`,
		listOwnerID, receivedAt,
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

// Вернуть обновленные товары из пошаренных на пользователя списков
// Т.к. это чужие списки, то возвращать только акцентованные по шарингу и неудаленные шаринги
func (s *ItemsReadRepository) GetUpdatedItemsForSharedListToUser(toUserId string, receivedAt int64) ([]models.ListItem, error) {
	var items []models.ListItem
	var statusAccepted = models.ShareStatusAccepted
	db := s.db

	rows, err := db.Query(
		`SELECT i.id, 
       			i.name, 
       			value,
       			is_marked, 
       			user_marked_id,
       			i.list_id, 
       			UNIX_TIMESTAMP(i.created_at), 
       			UNIX_TIMESTAMP(i.updated_at), 
       			UNIX_TIMESTAMP(i.received_at),
       			i.is_deleted
			FROM sl_item AS i 
			LEFT JOIN sl_shared_lists AS s  
			ON (i.list_id = s.list_id AND status = ? AND s.is_deleted = false) 
			WHERE s.to_user_id =? AND (i.received_at >= FROM_UNIXTIME(?) OR s.received_at >= FROM_UNIXTIME(?))`,
		statusAccepted, toUserId, receivedAt, receivedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repositories.ErrNotFound{}
		}

		return nil, err
	}
	defer rows.Close() // nolint errcheck

	items, err = s.scanItemRows(rows)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// Вернуть товары для списка
func (s *ItemsReadRepository) GetItemsForList(listID string) (*[]models.ListItem, error) {
	var items []models.ListItem
	db := s.db
	rows, err := db.Query(
		`SELECT id, 
       			name, 
       			value,
       			is_marked, 
       			user_marked_id,
       			list_id, 
       			UNIX_TIMESTAMP(created_at), 
       			UNIX_TIMESTAMP(updated_at), 
       			UNIX_TIMESTAMP(received_at),
       			is_deleted 
			FROM sl_item 
			WHERE list_id =?`,
		listID,
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

	return &items, nil
}

func (s *ItemsReadRepository) scanItemRows(rows *sql.Rows) ([]models.ListItem, error) {
	var items []models.ListItem

	for rows.Next() {
		var i models.ListItem
		err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Value,
			&i.IsMarked,
			&i.UserMarked,
			&i.ListID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ReceivedAt,
			&i.IsDeleted,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}
