package handlers

import (
	"ds-easy/src/database/repository"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s Service) RegisterExerciseRoutes() {
	baseUrl := "/exercises"

	s.Mux.HandleFunc(baseUrl, s.getExercisesHandler).Methods("GET")
	s.Mux.HandleFunc(baseUrl, s.addExerciseHandler).Methods("POST")
}

func (s Service) getExercisesHandler(w http.ResponseWriter, r *http.Request) {
	exercises, err := s.Queries.FindExercises(r.Context())
	if err != nil {
		log.Error("Errors occured", err)
	}

	jsonResp, err := json.Marshal(exercises)

	if err != nil {
		log.Error("Errors occured", err)
	}

	w.Write(jsonResp)
}

func (s Service) addExerciseHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("AddExerciseHandler")
	decoder := json.NewDecoder(r.Body)
	var payload repository.InsertExerciseParams
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Errors occured", err)
	}
	createdExercise, err := s.Queries.InsertExercise(r.Context(), payload)
	if err != nil {
		log.Error("Errors occured", err)
	}

	jsonResp, err := json.Marshal(createdExercise)

	if err != nil {
		log.Error("Errors occured", err)
	}

	w.Write(jsonResp)
}
