package handlers

import (
	"context"
	"ds-easy/src/database/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (s Service) RegisterUserRoutes() {
	baseUrl := "/users"

	s.Mux.HandleFunc(baseUrl, s.getUsersHandler).Methods("GET")
	s.Mux.HandleFunc(baseUrl, s.addUserHandler).Methods("POST")
	s.Mux.HandleFunc(baseUrl+"/{id}", s.getUserByIdHandler).Methods("GET")
}

func (s Service) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("GetUsersHandler")
	users, err := s.Queries.FindAllUsers(context.TODO())

	if err != nil {
		log.Error("Errors occured", err)
	}

	jsonResp, err := json.Marshal(users)

	if err != nil {
		log.Error("Errors occured", err)
	}

	w.Write(jsonResp)
}

func (s Service) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("GetUserByIdHandler")
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.Error("Errors occured", err)
	}

	user, err := s.Queries.FindUserById(context.TODO(), id)
	if err != nil {
		log.Error("Errors occured", err)
	}

	jsonResp, err := json.Marshal(user)

	if err != nil {
		log.Error("Errors occured", err)
	}

	w.Write(jsonResp)
}

func (s Service) addUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload repository.AddUserParams
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Errors occured", err)
	}
	createdUser, err := s.Queries.AddUser(r.Context(), payload)
	if err != nil {
		log.Error("Errors occured", err)
	}

	jsonResp, err := json.Marshal(createdUser)

	if err != nil {
		log.Error("Errors occured", err)
	}

	w.Write(jsonResp)
}
