package events

import (
	"fmt"
	"shopingList/pkg/models"
)

const (
	ShareListEventInvite     = "invite"
	ShareListEventAccept     = "accept"
	ShareListEventRefuse     = "refuse"
	ShareListEventDelete     = "share-delete"
	ShareListEventListDelete = "list-delete"
)

type ShareListEventType string

type ShareListEvent struct {
	typeEvent    ShareListEventType
	list         models.List
	user         models.User
	targetUserId string
}

func NewShareListEvent(typeEvent ShareListEventType, list models.List, user models.User, targetUserId string) ShareListEvent {
	return ShareListEvent{typeEvent: typeEvent, list: list, user: user, targetUserId: targetUserId}
}

func (s *ShareListEvent) GetEventType() string {
	return fmt.Sprint("%s")
}

func (s *ShareListEvent) TypeEvent() ShareListEventType {
	return s.typeEvent
}

func (s *ShareListEvent) List() models.List {
	return s.list
}

func (s *ShareListEvent) User() models.User {
	return s.user
}

func (s *ShareListEvent) TargetUserId() string {
	return s.targetUserId
}
