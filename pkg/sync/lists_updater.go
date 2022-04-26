package sync

import (
	"errors"
	"shopingList/pkg/models"
	"shopingList/pkg/repositories"
	"time"
)

type ListsUpdater struct {
	userId          string
	listsRepository *repositories.ListsRepository
	listsCollection *ListsCollection
}

func NewUpdaterList(userId string, listsRepository *repositories.ListsRepository, listsCollection *ListsCollection) ListsUpdater {
	return ListsUpdater{userId: userId, listsRepository: listsRepository, listsCollection: listsCollection}
}

func (s *ListsUpdater) Run(lists []models.List) error {
	if len(lists) == 0 {
		return nil
	}

	var listsUpdate []models.List
	var listsCreate []models.List

	for _, list := range lists {
		if list.OwnerID != s.userId {
			return errors.New("forbidden to update another user's list")
		}

		list.ReceivedAt = time.Now().UTC().Unix()
		existList, err := s.listsCollection.GetListForId(list.ID, s.userId)
		if err != nil {
			if _, ok := err.(repositories.ErrNotFound); ok {
				listsCreate = append(listsCreate, list)
				continue
			}

			return errors.New("can't get list" + err.Error())
		}

		if list.IsEqual(*existList) {
			continue
		}

		if existList.IsTemplate != list.IsTemplate {
			return errors.New("forbidden change is_template value for list: " + existList.ID)
		}

		listsUpdate = append(listsUpdate, list)
	}

	for _, list := range listsUpdate {
		err := s.listsRepository.UpdateList(&list)

		if err != nil {
			return errors.New("Error update list; " + err.Error())
		}

		s.listsCollection.AddList(&list)
	}

	for _, list := range listsCreate {
		err := s.listsRepository.CreateList(&list)
		if err != nil {
			return errors.New("Error create list; " + err.Error())
		}

		s.listsCollection.AddList(&list)
	}

	return nil
}
