package sync

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"shopingList/api"
	"shopingList/api/auth"
	"shopingList/api/controllers"
	"shopingList/pkg/events"
	"shopingList/pkg/sync"
	"shopingList/store"
)

type SyncController struct {
	authService     *auth.Service
	dataService     store.DataService
	chanGoodsChange chan events.GoodsChangeEvent
	chanShareChange chan events.ShareListEvent
}

func NewSyncController(
	authService *auth.Service,
	dataService store.DataService,
	chanGoodsChange chan events.GoodsChangeEvent,
	chanShareChange chan events.ShareListEvent) *SyncController {
	return &SyncController{
		authService:     authService,
		dataService:     dataService,
		chanGoodsChange: chanGoodsChange,
		chanShareChange: chanShareChange}
}

// Routes returns slice of server routes
func (s *SyncController) Routes() []api.Route {
	return []api.Route{
		{
			Name:   "GetShoppingListUpdates",
			Method: "GET",
			Path:   "/shoppingList/updates",
			Func:   s.getSyncUpdates,
		},
		{
			Name:   "PostShoppingListUpdates",
			Method: "POST",
			Path:   "/shoppingList/updates",
			Func:   s.saveSyncUpdates,
		},
	}
}

func (s *SyncController) getSyncUpdates(w http.ResponseWriter, r *http.Request) {
	currentUser, err := controllers.GetAuthorizedUser(s.authService, r)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusUnauthorized, err, "can't auth user by token", api.ErrUserNotFound)
		return
	}

	var query ShoppingListUpdatesRequest
	err = schema.NewDecoder().Decode(&query, r.URL.Query())
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't query params", api.ErrDecode)
		return
	}

	var syncReceiver sync.Receiver
	pack, err := syncReceiver.GetUpdates(s.dataService, *currentUser, query.Date)
	if err != nil {
		log.Errorln(errors.Wrap(err, "Error in getSyncUpdates()"))
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "error getting updates", api.ErrInternal)
		return
	}

	api.SendDataJSON(w, r, http.StatusOK, pack)
}

func (s *SyncController) saveSyncUpdates(w http.ResponseWriter, r *http.Request) {
	currentUser, err := controllers.GetAuthorizedUser(s.authService, r)
	if err != nil {
		api.SendErrorJSON(w, r, http.StatusUnauthorized, err, "can't auth user by token", api.ErrUserNotFound)
		return
	}

	var data ShoppingListUpdates
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		api.SendErrorJSON(w, r, http.StatusBadRequest, err, "error decode request", api.ErrDecode)
		return
	}

	result, errs := data.Validate()
	if !result {
		api.SendErrorJSON(w, r, http.StatusBadRequest, errors.New(errs[0]), "validation error", api.ErrValidationData)
		return
	}

	syncUpdater := sync.NewUpdater(s.dataService, *currentUser)
	syncUpdater.ChanGoodsChange = s.chanGoodsChange
	syncUpdater.ChanShareChange = s.chanShareChange
	err = syncUpdater.RunUpdate(data.Users, data.Lists, data.Shares, data.Items, data.UserProducts)

	if err != nil {
		log.Errorln(errors.Wrap(err, "Error in saveSyncUpdates()"))
		api.SendErrorJSON(w, r, http.StatusInternalServerError, err, "error synchronization", api.ErrInternal)
		return
	}

	api.SendDataJSON(w, r, http.StatusOK, nil)
}
