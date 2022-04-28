package controllers

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"net/http"
	"shopingList/api"
	"shopingList/api/auth"
	"shopingList/pkg/events"
	"shopingList/pkg/models"
	"shopingList/pkg/readModels"
	"shopingList/pkg/repositories"
	"shopingList/store"
	"time"
)

type SharedListsController struct {
	authService          *auth.Service
	sharesRepository     repositories.SharesRepository
	sharesReadRepository readModels.SharesReadRepository
	itemsReadRepository  readModels.ItemsReadRepository
	listsRepository      readModels.ListsReadRepository
	ChanShareChange      chan events.ShareListEvent
}

func NewSharedListsController(authService *auth.Service, dataService store.DataService) *SharedListsController {
	return &SharedListsController{authService: authService,
		sharesRepository:     dataService.GetSharesRepository(nil),
		sharesReadRepository: dataService.GetSharesReadRepository(),
		itemsReadRepository:  dataService.GetItemsReadRepository(),
		listsRepository:      dataService.GetListsReadRepository()}
}

func (s *SharedListsController) Routes() []api.Route {
	return []api.Route{
		{
			Name:   "Status",
			Method: "POST",
			Path:   "/share-list/{share_id}/accept",
			Func:   s.accept,
		},
	}
}

func (s *SharedListsController) accept(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	share_id := vars["share_id"]

	err := validation.Validate(share_id, validation.Required, validation.Length(36, 36))
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "wrong format share list", api.ErrDecode)
		return
	}

	currentUser, err := GetAuthorizedUser(s.authService, r)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusUnauthorized, err, "can't auth user by token", api.ErrUserNotFound)
		return
	}

	share, err := s.sharesReadRepository.GetShareForUser(share_id, currentUser.ID)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, err.Error(), api.ErrInternal)
		return
	}

	if share.IsDeleted {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "The share is deleted", api.ErrNoPermission)
		return
	}

	share.Status = models.ShareStatusAccepted
	share.UpdatedAt = time.Now().UTC().Unix()

	err = s.sharesRepository.UpdateShare(share)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "error saving share", api.ErrInternal)
		return
	}

	list, err := s.listsRepository.GetListForIdAndOwner(share.ListID, share.OwnerID)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "error get list", api.ErrInternal)
		return
	}

	items, err := s.itemsReadRepository.GetItemsForList(share.ListID)
	if err != nil {
		if _, ok := err.(repositories.ErrNotFound); ok {
			api.SendDataJSON(w, r, http.StatusOK, nil)
			return
		}

		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get items for share object", api.ErrInternal)
		return
	}

	api.SendDataJSON(w, r, http.StatusOK, map[string]*[]models.ListItem{"items": items})

	if s.ChanShareChange != nil {
		event := events.NewShareListEvent(events.ShareListEventAccept, list, *currentUser, share.OwnerID)
		s.ChanShareChange <- event
	}

	return
}
