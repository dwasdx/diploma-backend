package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"shopingList/pkg"
	"shopingList/pkg/forms"
	"shopingList/pkg/models"
)

type UsersRepository struct {
	db models.DB
}

func NewUsersRepository(db models.DB) UsersRepository {
	if db == nil {
		panic("db param is nil")
	}

	return UsersRepository{db: db}
}

func (s *UsersRepository) GetUser(id string) (*models.User, error) {
	row := s.db.QueryRow(
		`SELECT id, 
       			name, 
       			email, 
       			phone, 
       			code, 
       			UNIX_TIMESTAMP(created_at), 
       			UNIX_TIMESTAMP(updated_at), 
       			is_activated, 
       			is_deleted 
        FROM sl_users
        WHERE id=?`, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Phone,
		&u.Code,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.IsActivated,
		&u.IsDeleted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.ErrNotFoundInStorage
		}

		return nil, err
	}
	return &u, nil
}

// CreateUser creates new user in database by getting info form input model and return created user model.
func (s *UsersRepository) CreateUser(f forms.RegForm) (*models.User, error) {
	db := s.db

	uuidValue := uuid.New()

	_, err := db.Exec(
		`INSERT sl_users 
			SET id=?, 
			    name=?,
			    email=?, 
			    phone=?, 
			    code=?`,
		uuidValue,
		f.Username,
		f.Email,
		f.Phone,
		nil,
	)
	if err != nil {
		return nil, err
	}

	user, err := s.GetUser(uuidValue.String())
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser update user in database from input user model
func (s *UsersRepository) UpdateUser(u *models.User) error {
	_, err := s.db.Exec(
		`UPDATE sl_users 
   		SET name=?, phone=?, code=?, is_activated=?, updated_at=FROM_UNIXTIME(?), is_deleted=?  
    	WHERE id=?`,
		u.Name, u.Phone, u.Code.SqlValue(), u.IsActivated, u.UpdatedAt, u.IsDeleted, u.ID)
	if err != nil {
		return err
	}
	return nil
}
