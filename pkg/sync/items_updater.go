package sync

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"shopingList/pkg/models"
	"shopingList/pkg/readModels"
	"shopingList/pkg/repositories"
	"time"
)

type ItemsUpdater struct {
	user                 models.User
	userIdsForList       map[string][]string
	listsCollection      *ListsCollection
	itemsRepository      repositories.ItemsRepository
	listsReadRepository  readModels.ListsReadRepository
	sharesReadRepository readModels.SharesReadRepository
	usersReadRepository  readModels.UsersReadRepository
	notificationService  NotificationsCreateService
	eventCollection      *EventCollection
}

func NewItemsUpdater(
	user models.User,
	listsCollection *ListsCollection,
	itemsRepository repositories.ItemsRepository,
	listsReadRepository readModels.ListsReadRepository,
	sharesReadRepository readModels.SharesReadRepository,
	usersReadRepository readModels.UsersReadRepository,
	eventCollection *EventCollection) *ItemsUpdater {
	return &ItemsUpdater{
		user:                 user,
		userIdsForList:       make(map[string][]string),
		listsCollection:      listsCollection,
		itemsRepository:      itemsRepository,
		listsReadRepository:  listsReadRepository,
		sharesReadRepository: sharesReadRepository,
		usersReadRepository:  usersReadRepository,
		eventCollection:      eventCollection,
		notificationService:  NotificationsCreateService{}}
}

func (s *ItemsUpdater) getUserId() string {
	return s.user.ID
}

// Обработать товары
// Может создавать и обновлять только товары для собственных списков
// и списков, которые расшарены для юзера и активны (шаринги не удалены)
func (s *ItemsUpdater) Run(items []models.ListItem, lists []models.List) error {
	if len(items) == 0 {
		return nil
	}

	var listIdsFormItems []string

	for _, item := range items {
		listIdsFormItems = append(listIdsFormItems, item.ListID)
	}

	listOwnIds := make(map[string]bool)
	mapLists := make(map[string]models.List)

	ownLists, err := s.listsReadRepository.GetListsForIdsAndOwner(listIdsFormItems, s.getUserId())
	if err != nil {
		return errors.New("Error get own lists for items; " + err.Error())
	}
	for _, list := range ownLists {
		listOwnIds[list.ID] = true
		mapLists[list.ID] = list
	}

	// Добавляем списки, которые пришли в пакете синхронизации
	for _, list := range lists {
		if list.OwnerID == s.getUserId() {
			listOwnIds[list.ID] = true
			mapLists[list.ID] = list
		}
	}

	sharesMap := make(map[string][]models.ListShare)
	shares, err := s.sharesReadRepository.GetSharesForUserForListIds(listIdsFormItems, s.getUserId())
	if err != nil {
		return errors.New("Error get shares for the item`s list; " + err.Error())
	}

	for _, share := range shares {
		sharesMap[share.ListID] = append(sharesMap[share.ListID], share)
	}

	// Делим товары на 2 массива
	var itemsForOwnList []models.ListItem
	var ItemsForSharedLists []models.ListItem

	// Проверить, права доступа на обновление товаров
	for _, item := range items {
		if _, ok := listOwnIds[item.ListID]; ok {
			// Нельзя отметить товар купленным для списка-шаблона
			if item.IsMarked {
				if tmpList, ok := mapLists[item.ListID]; ok {
					if tmpList.IsTemplate {
						return errors.New(
							fmt.Sprintf("forbidden to mark a product for a template list. Item id: %s, List id: %s",
								item.ID, tmpList.ID))
					}
				}
			}

			itemsForOwnList = append(itemsForOwnList, item)
			continue
		}

		if shareList, ok := sharesMap[item.ListID]; ok {
			for _, share := range shareList {
				if share.IsDeleted {
					continue
				}

				// Принимаем для акцептованных шарингов
				// И для отказанных, т.к. отказ мог быть сделан после изменения товаров и их надо принять
				if share.Status == models.ShareStatusNew {
					continue
				}

				ItemsForSharedLists = append(ItemsForSharedLists, item)
			}

			continue
		}

		// ID cписка из товара нет ни в собственных листах, ни в расшаренных
		return errors.New(fmt.Sprintf("Error sync item with id: %s. Its list with id: %s doesn`t found.", item.ID, item.ListID))
	}

	// Проверить id-пользователей, указанных в товаре в user_marked
	userMarkedIds := make([]string, 0)
	userMarkedExists := make(map[string]bool)

	for _, item := range items {
		if !item.UserMarked.IsEmpty() {
			userMarkedIds = append(userMarkedIds, item.UserMarked.String)
			userMarkedExists[item.UserMarked.String] = false
		}
	}

	if len(userMarkedIds) > 0 {
		markedUsers, err := s.usersReadRepository.GetUsersForIds(userMarkedIds...)
		if err != nil {
			if _, ok := err.(repositories.ErrNotFound); !ok {
				return errors.New("error get user_marked by ids; " + err.Error())
			}
		} else {
			for _, user := range markedUsers {
				userMarkedExists[user.ID] = true
			}

			for userId, exists := range userMarkedExists {
				if !exists {
					return errors.New("user_marked doesn`t found by id; " + userId)
				}
			}
		}
	}

	if len(itemsForOwnList) > 0 {
		err = s.updateOwnItems(&itemsForOwnList)
		if err != nil {
			return errors.New("Error with update items in own lists; " + err.Error())
		}
	}

	if len(ItemsForSharedLists) > 0 {
		err = s.updateItemsFromSharedList(&ItemsForSharedLists)
		if err != nil {
			return errors.New("Error with update items in shared lists; " + err.Error())
		}
	}

	return nil
}

// Обновить товары из собственных списков
func (s *ItemsUpdater) updateOwnItems(items *[]models.ListItem) error {
	for _, item := range *items {
		list, err := s.listsCollection.GetListForId(item.ListID, s.getUserId())
		item.ReceivedAt = time.Now().UTC().Unix()

		existItem, err := s.itemsRepository.GetItem(item.ID, item.ListID)
		if err != nil {
			if _, ok := err.(repositories.ErrNotFound); !ok {
				return errors.New("can`t get item for id: " + item.ID)
			}

			err = s.itemsRepository.CreateItem(&item)
			if err != nil {
				return errors.New("Can`t create item; " + err.Error())
			}

			s.createNotificationForItem(&item, nil, list)
			continue
		}

		if existItem.IsEqual(&item) {
			continue
		}

		if list.IsTemplate && item.IsMarked {
			return errors.New(
				"Forbidden to mark products belonging to the template. Item: " + item.ID + ", Template list: " + list.ID)
		}

		err = s.itemsRepository.UpdateItem(&item)
		if err != nil {
			return errors.New("Can`t update item; " + err.Error())
		}

		s.createNotificationForItem(&item, &existItem, list)
	}
	return nil
}

// Обновить товары из пошаренных списков
func (s *ItemsUpdater) updateItemsFromSharedList(items *[]models.ListItem) error {
	for _, item := range *items {
		item.ReceivedAt = time.Now().UTC().Unix()
		lists, err := s.listsReadRepository.GetListsSharedForUserForIds([]string{item.ListID}, s.getUserId())
		if err != nil {
			return err
		}

		if len(lists) == 0 {
			return errors.New("Shared list doesn`t found for id: " + item.ListID)
		}

		existItem, err := s.itemsRepository.GetItem(item.ID, item.ListID)
		if err != nil {
			if _, ok := err.(repositories.ErrNotFound); !ok {
				return errors.New("can`t get item for id: " + item.ID)
			}

			// Разрешаем создавать товары в пошаренных списках
			err = s.itemsRepository.CreateItem(&item)
			if err != nil {
				return errors.New("Can`t create item in shared list; " + err.Error())
			}

			s.createNotificationForItem(&item, nil, &lists[0])
			continue
		}

		if existItem.IsEqual(&item) {
			continue
		}

		// Для товаров из пошаренных списков пока разрешаем менять все
		err = s.itemsRepository.UpdateItem(&item)
		if err != nil {
			return errors.New("Can`t update item in shared list; " + err.Error())
		}

		s.createNotificationForItem(&item, &existItem, &lists[0])
	}

	return nil
}

func (s *ItemsUpdater) createNotificationForItem(item *models.ListItem, existItem *models.ListItem, list *models.List) {
	listName := s.listsCollection.GetListNameById(item.ListID, list.OwnerID)

	// Получатели - акцептованные пользователя списка
	targetIds := s.getUserIdsForList(list.ID, list.OwnerID)

	// Если список чужой, то добавить и владельца
	if list.OwnerID != s.getUserId() {
		targetIds = append(targetIds, list.OwnerID)
	}

	event := s.notificationService.createNotificationForItem(item, existItem, listName, &s.user, targetIds)
	s.eventCollection.AddGoodEvent(*event)
}

// Вернуть пользователей пошаренного и акцептованного списка.
// Текущий пользователь не должен входить в список
func (s *ItemsUpdater) getUserIdsForList(listId string, ownerId string) []string {
	_, ok := s.userIdsForList[listId]

	ids := make([]string, 0)

	if !ok {
		rawIds, err := s.sharesReadRepository.GetAcceptedUserIdsFromSharedList(listId, ownerId)

		if err != nil {
			log.Errorln(err)
			return nil
		}

		for _, id := range rawIds {
			if id != s.user.ID {
				ids = append(ids, id)
			}
		}

		s.userIdsForList[listId] = ids
	}

	return s.userIdsForList[listId]
}
