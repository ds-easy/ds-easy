package handlers

import (
	"context"
	"ds-easy/src/database/repository"
	utils "ds-easy/src/web/handlers/util"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (s Service) registerUserRoutes() {
	baseUrl := "/users"

	s.Mux.HandleFunc(baseUrl, s.getUsersHandler).Methods("GET")
	s.Mux.HandleFunc(baseUrl, s.addUserHandler).Methods("POST")
	s.Mux.HandleFunc(baseUrl+"/{id}", s.getUserByIdHandler).Methods("GET")
}

func (s Service) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("GetUsersHandler")

	users, err := s.Queries.FindAllUsers(context.TODO())
	if err != nil {
		log.Error("Error getting users from DB: ", err)
	}

	jsonResp, err := json.Marshal(users)
	if err != nil {
		log.Error("Error serializing users: ", err)
	}

	w.Write(jsonResp)
}

func (s Service) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("GetUserByIdHandler")
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.Error("Error getting query parameter 'id': ", err)
	}

	user, err := s.Queries.FindUserById(context.TODO(), id)
	if err != nil {
		log.Error("Error getting user from DB: ", err)
	}

	jsonResp, err := json.Marshal(user)
	if err != nil {
		log.Error("Error serializing user instance: ", err)
	}

	w.Write(jsonResp)
}

func (s Service) addUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("AddUserHandler")
	decoder := json.NewDecoder(r.Body)

	var pbAddUserParams struct {
		repository.AddUserParams
		Password string `json:"password"`
	}

	err := decoder.Decode(&pbAddUserParams)
	if err != nil {
		log.Error("Error decoding request body: ", err)
	}

	userAddParams := repository.AddUserParams{
		PbID:      "",
		FirstName: pbAddUserParams.FirstName,
		LastName:  pbAddUserParams.LastName,
		Email:     pbAddUserParams.Email,
		Admin:     0,
	}

	userAddParams.PbID, err = utils.PBAddUser(userAddParams, pbAddUserParams.Password)
	if err != nil {
		log.Error("Error creating user in PB: ", err)
	}

	createdUser, err := s.Queries.AddUser(r.Context(), userAddParams)
	if err != nil {
		log.Error("Error adding user to DB: ", err)
	}

	jsonResp, err := json.Marshal(createdUser)
	if err != nil {
		log.Error("Error serializing newly created user: ", err)
	}

	w.Write(jsonResp)
}
