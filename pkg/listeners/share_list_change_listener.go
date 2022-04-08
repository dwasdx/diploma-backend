package listeners

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"shopingList/pkg/events"
	"shopingList/pkg/models"
	"shopingList/pkg/repositories"
	"shopingList/pkg/services"
)

type ShareListChangeListener struct {
	Repository  repositories.NotificationsRepository
	PushChannel chan services.PushNotificationMessage
}

func (s *ShareListChangeListener) Run(channel chan events.ShareListEvent) {
	for event := range channel {
		log.Info("Receive message for ShareListChangeListener")

		err := s.Handle(&event)
		if err != nil {
			log.Error("Error handle event", err)
		}
	}
}

func (s *ShareListChangeListener) Handle(event interface{}) error {
	model, ok := event.(*events.ShareListEvent)

	if !ok {
		err := errors.New(fmt.Sprintf("event must be a type not ShareListEvent, event is type: %#v", event))
		log.Fatalf("%+v", err)
	}

	typeEvent := model.TypeEvent()
	list := model.List()
	user := model.User()
	targetUserId := model.TargetUserId()

	var message string
	var typeNotification models.NotificationType

	switch typeEvent {
	case events.ShareListEventInvite:
		message = fmt.Sprintf("Пользователь %d пригласил вас в список \"%s\"", user.Phone, list.Name)
		typeNotification = models.NotificationTypeListJoining
	case events.ShareListEventAccept:
		message = fmt.Sprintf("Пользователь %d присоединился к списку \"%s\"", user.Phone, list.Name)
		typeNotification = models.NotificationTypeListJoining
	case events.ShareListEventRefuse:
		message = fmt.Sprintf("Пользователь %d покинул список \"%s\"", user.Phone, list.Name)
		typeNotification = models.NotificationTypeListDetachment
	case events.ShareListEventDelete:
		message = fmt.Sprintf("Пользователь %d отменил шаринг списка \"%s\"", user.Phone, list.Name)
		typeNotification = models.NotificationTypeListShareDelete
	case events.ShareListEventListDelete:
		message = fmt.Sprintf("Пользователь %d удалил список \"%s\"", user.Phone, list.Name)
		typeNotification = models.NotificationTypeListDelete
	default:
		err := errors.New(fmt.Sprintf("typeEvent is wrong, type: %s", typeEvent))
		log.Fatalf("%+v", err)
	}

	form := models.NotificationCreateForm{
		TypeNotification: typeNotification,
		Message:          message,
		UserId:           user.ID,
		UserPhone:        user.Phone,
		ListId:           list.ID,
		TargetUserId:     targetUserId,
	}

	err := s.Repository.Create(&form)

	if err != nil {
		return errors.Wrap(err, "Error create notification in the ShareListChangeListener")
	}

	if s.PushChannel == nil {
		log.Warn("[ERROR] PushChannel in ShareListChangeListener is nil")
	} else {
		pushMessage := services.PushNotificationMessage{Notification: form, TargetUserIds: []string{targetUserId}}
		s.PushChannel <- pushMessage
	}

	return nil
}
