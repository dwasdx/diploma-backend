package sync

import (
	"shopingList/pkg/models"
	"shopingList/pkg/readModels"
)

type ListsCollection struct {
	lists      map[string]*models.List
	repository *readModels.ListsReadRepository
}

func NewListsCollection(repository *readModels.ListsReadRepository) *ListsCollection {
	return &ListsCollection{
		lists:      make(map[string]*models.List),
		repository: repository}
}

func (s *ListsCollection) GetListForId(listId string, ownerId string) (*models.List, error) {
	if s.lists == nil {
		s.lists = make(map[string]*models.List)
	}

	_, ok := s.lists[listId]

	if !ok {
		existList, err := s.repository.GetListForIdAndOwner(listId, ownerId)

		if err != nil {
			return nil, err
		}

		s.lists[listId] = &existList
	}

	return s.lists[listId], nil
}

func (s *ListsCollection) GetListNameById(listId string, ownerId string) string {
	list, err := s.GetListForId(listId, ownerId)

	if err != nil {
		return ""
	} else {
		return list.Name
	}
}

func (s *ListsCollection) AddList(list *models.List) {
	s.lists[list.ID] = list
}
