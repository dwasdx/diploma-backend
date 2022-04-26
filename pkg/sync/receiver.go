package sync

import (
	"github.com/pkg/errors"
	"shopingList/pkg/models"
	"shopingList/pkg/repositories"
	"shopingList/store"
)

type Receiver struct {
	dataService store.DataService
}

func (s *Receiver) GetUpdates(dataService store.DataService, user models.User, updatedAt int64) (*UpdatesPack, error) {
	var resp UpdatesPack
	s.dataService = dataService

	// Шаринги (свои и чужие) для пользователя
	shares, err := s.getUpdatedSharingsForUser(user.ID, updatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting sharing objects in receiver")
	}
	resp.Shares = append(resp.Shares, shares...)

	// Списки (свои и чужие)
	lists, err := s.getUpdatedListForUser(user.ID, updatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting list objects in receiver")
	}
	resp.Lists = append(resp.Lists, lists...)

	// Элементы (свои и чужие)
	items, err := s.getUpdatedItemsForUser(user.ID, updatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting items objects in receiver")
	}
	resp.Items = append(resp.Items, items...)

	// Избранное
	userProducts, err := s.getUpdatedUserProductsForUser(user.ID, updatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting userProducts objects in receiver")
	}
	resp.UserProducts = append(resp.UserProducts, userProducts...)

	itemIds := make(map[string]string)

	//////////////////////////////////////////////////////////////
	// Добавить в выдачу объекты списков для измененных элементов
	//////////////////////////////////////////////////////////////
	for _, item := range resp.Items {
		itemIds[item.ListID] = item.ListID
	}

	var listForAdd []string
	for _, listId := range itemIds {
		if !resp.IsExistList(listId) {
			listForAdd = append(listForAdd, listId)
		}
	}

	listReadRepository := s.dataService.GetListsReadRepository()

	listForItems, err := listReadRepository.GetListsForIdsAndOwner(listForAdd, user.ID)
	resp.Lists = append(resp.Lists, listForItems...)

	usersRepository := s.dataService.GetUsersReadRepository()

	// Получить id-ы пользователей для всех объектов resp
	userIds := resp.GetUserIdsInObjects()
	users, err := usersRepository.GetUsersForIds(userIds...)

	if err != nil {
		return nil, errors.Wrap(err, "can't get users in receiver")
	}

	// Сохранить в resp разный объем данных для объектов юзеров
	for _, u := range users {
		if u.ID == user.ID {
			resp.Users = append(resp.Users, u)
		} else {
			resp.Users = append(resp.Users, models.CreateExternalUser(&u))
		}
	}

	return &resp, nil
}

// Получить шаринги для пользователя
func (s *Receiver) getUpdatedSharingsForUser(userId string, updatedAt int64) ([]models.ListShare, error) {
	var shares []models.ListShare

	repository := s.dataService.GetSharesReadRepository()

	// Собственные шаринги
	myShares, err := repository.GetUpdatedSharesForOwner(userId, updatedAt)
	if err != nil {
		if _, ok := err.(repositories.ErrNotFound); !ok {
			return nil, err
		}
	}
	shares = append(shares, myShares...)

	// Шаринги на текущего пользователя
	sharesToUser, err := repository.GetUpdatedSharesToUser(userId, updatedAt)
	if err != nil {
		if _, ok := err.(repositories.ErrNotFound); !ok {
			return nil, err
		}
	}

	shares = append(shares, sharesToUser...)
	return shares, nil
}

// Получить списки для пользователя
func (s *Receiver) getUpdatedListForUser(userId string, updatedAt int64) ([]models.List, error) {
	var lists []models.List

	listReadRepository := s.dataService.GetListsReadRepository()

	// Получить списки пользователя
	selfLists, err := listReadRepository.GetUpdatedListsForOwner(userId, updatedAt) // my lists
	if err != nil {
		if _, ok := err.(repositories.ErrNotFound); !ok {
			return nil, err
		}
	}

	lists = append(lists, selfLists...)

	// Получить списки, пошаренные на пользователя
	sharedLists, err := listReadRepository.GetUpdatedListsSharedToUser(userId, updatedAt)
	if err != nil {
		if _, ok := err.(repositories.ErrNotFound); !ok {
			return nil, err
		}
	}

	lists = append(lists, sharedLists...)

	return lists, nil
}

// Получить элементы для пользователя
func (s *Receiver) getUpdatedItemsForUser(userId string, updatedAt int64) ([]models.ListItem, error) {
	items := make(map[string]models.ListItem, 0)

	repository := s.dataService.GetItemsReadRepository()

	// Измененные элементы своих списков
	ownItems, err := repository.GetUpdatedItemsForUser(userId, updatedAt)
	if err != nil {
		if _, ok := err.(repositories.ErrNotFound); !ok {
			return nil, err
		}
	}

	for _, item := range ownItems {
		items[item.ID] = item
	}

	// Измененные элементы пошаренных списков
	sharedItems, err := repository.GetUpdatedItemsForSharedListToUser(userId, updatedAt)
	if err != nil {
		if _, ok := err.(repositories.ErrNotFound); !ok {
			return nil, err
		}
	}

	for _, item := range sharedItems {
		if _, ok := items[item.ID]; !ok {
			items[item.ID] = item
		}
	}

	var result []models.ListItem
	for _, item := range items {
		result = append(result, item)
	}

	return result, nil
}

// Получить локальные продукты пользователя
func (s *Receiver) getUpdatedUserProductsForUser(userId string, updatedAt int64) ([]models.UserProduct, error) {
	var userProduct []models.UserProduct
	userProductReadRepository := s.dataService.UserProductsReadRepository()

	userProduct, err := userProductReadRepository.GetUpdatedUserProductsForUser(userId, updatedAt) // my userProducts
	if err != nil {
		if _, ok := err.(repositories.ErrNotFound); !ok {
			return nil, err
		}
	}

	return userProduct, nil
}
