package controllers

import (
	"net/http"
	"shopingList/api"
	"time"
)

// PublicController API-service
type PublicController struct {
}

// NewPublic returns new instance of PublicController API
func NewPublic() *PublicController {
	return &PublicController{}
}

// Routes returns slice of server routes
func (s *PublicController) Routes() []api.Route {
	return []api.Route{
		{
			Name:   "Ping",
			Method: "GET",
			Path:   "/ping",
			Func:   s.ping,
		},
	}
}

func (s *PublicController) ping(w http.ResponseWriter, r *http.Request) {
	api.SendJSON(w, r, api.JSON{"pong": time.Now().String()})
}
