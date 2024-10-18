package userHandlers

import (
	"ds-easy/src/database/repository"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	queries repository.Queries
	mux     mux.Router
}

func (s userService) RegisterRoutes() {
	baseUrl := "users"

	s.mux.HandleFunc(baseUrl, s.getUsersHandler)
}

func (s userService) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("GetUsersHandler")
	users, err := s.queries.FindAllUsers(nil)

	if err != nil {
		log.Error("Errors occured", err)
	}

	jsonResp, err := json.Marshal(users)

	if err != nil {
		log.Error("Errors occured", err)
	}

	w.Write(jsonResp)
}
