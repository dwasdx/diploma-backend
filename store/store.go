package store

import (
	"database/sql"
	"shopingList/pkg/readModels"
	"shopingList/pkg/repositories"
)

// DataService required data service interface
type DataService interface {
	CreateTransaction() (*sql.Tx, error)

	// Репозитории
	GetListsRepository(tx *sql.Tx) repositories.ListsRepository
	GetItemsRepository(tx *sql.Tx) repositories.ItemsRepository
	GetSharesRepository(tx *sql.Tx) repositories.SharesRepository
	GetUsersRepository(tx *sql.Tx) repositories.UsersRepository
	UserProductsRepository(tx *sql.Tx) repositories.UserProductsRepository

	// Репозитории на чтении
	GetListsReadRepository() readModels.ListsReadRepository
	GetItemsReadRepository() readModels.ItemsReadRepository
	GetSharesReadRepository() readModels.SharesReadRepository
	GetUsersReadRepository() readModels.UsersReadRepository
	UserProductsReadRepository() readModels.UserProductsReadRepository
}
