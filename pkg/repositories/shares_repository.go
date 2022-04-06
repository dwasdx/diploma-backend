package repositories

import (
	"shopingList/pkg/models"
)

type SharesRepository struct {
	db models.DB
}

func NewSharesRepository(db models.DB) SharesRepository {
	if db == nil {
		panic("db param is nil")
	}

	return SharesRepository{db: db}
}

func (s *SharesRepository) CreateShare(share *models.ListShare) error {
	_, err := s.db.Exec(
		`INSERT INTO sl_shared_lists (
                    id, list_id, to_user_id, owner_id, status, created_at, updated_at, received_at, is_deleted
                    )
        VALUES (?, ?, ?, ?, ?, FROM_UNIXTIME(?), FROM_UNIXTIME(?), FROM_UNIXTIME(?), ?)`,
		share.ID, share.ListID, share.ToUserID, share.OwnerID, share.Status,
		share.CreatedAt, share.UpdatedAt, share.ReceivedAt, share.IsDeleted)
	if err != nil {
		return err
	}

	return nil
}

func (s *SharesRepository) UpdateShare(share *models.ListShare) error {
	_, err := s.db.Exec(
		`UPDATE sl_shared_lists 
		SET status=?, updated_at=FROM_UNIXTIME(?), received_at=FROM_UNIXTIME(?), is_deleted=? 
		WHERE id=? AND owner_id=?`,
		share.Status, share.UpdatedAt, share.ReceivedAt, share.IsDeleted, share.ID, share.OwnerID)
	if err != nil {
		return err
	}

	return nil
}
