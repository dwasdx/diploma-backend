// Репозитарий получения объектов models.User
package readModels

import (
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/pkg/errors"
	"shopingList/pkg"
	"shopingList/pkg/models"
)

// Имя таблицы в реляционной БД
const TableUsers = "sl_users"

type UsersReadRepository struct {
	db *sql.DB
}

// Создать новый репозитарий
func NewUsersReadRepository(db *sql.DB) UsersReadRepository {
	return UsersReadRepository{db: db}
}

// Вернуть объект User по номеру телефона
func (s *UsersReadRepository) GetUserByPhone(phone int64) (*models.User, error) {
	var user models.User

	selDataset := s.getSelectDataset()
	_, err := selDataset.Where(goqu.Ex{"phone": phone}).ScanStruct(&user)
	if err != nil {
		return nil, errors.Wrap(err, "Error search users by id")
	}

	if user.ID == "" {
		return nil, pkg.ErrNotFoundInStorage
	}

	return &user, nil
}

// Вернуть объект User по id
func (s *UsersReadRepository) GetUser(id string) (*models.User, error) {
	var user models.User

	selDataset := s.getSelectDataset()

	_, err := selDataset.Where(goqu.Ex{"id": id}).ScanStruct(&user)
	if err != nil {
		return nil, errors.Wrap(err, "Error search users by id")
	}

	return &user, nil
}

// Вернуть список юзеров по массиву ID-ов
func (s *UsersReadRepository) GetUsersForIds(UserIds ...string) ([]models.User, error) {
	var users []models.User

	if len(UserIds) == 0 {
		return users, nil
	}

	selDataset := s.getSelectDataset()
	err := selDataset.Where(goqu.Ex{"id": UserIds}).ScanStructs(&users)
	if err != nil {
		return nil, errors.Wrap(err, "Error search users by ids")
	}

	return users, nil
}

// Вернуть список юзеров по массиву телефонов
func (s *UsersReadRepository) GetUsersByPhones(UserPhones ...int64) ([]models.User, error) {
	var users []models.User

	if len(UserPhones) == 0 {
		return users, nil
	}

	selDataset := s.getSelectDataset()
	err := selDataset.Where(goqu.Ex{"phone": UserPhones}).ScanStructs(&users)
	if err != nil {
		return nil, errors.Wrap(err, "Error search users by phones")
	}

	return users, nil
}

func (s *UsersReadRepository) getSelectDataset() *goqu.SelectDataset {
	db := goqu.New("mysql", s.db)

	return db.From(TableUsers).Select("id", "name", "email", "phone", "code",
		goqu.L("UNIX_TIMESTAMP(`created_at`)").As("created_at"),
		goqu.L("UNIX_TIMESTAMP(`updated_at`)").As("updated_at"),
		"is_activated", "is_deleted")
}
