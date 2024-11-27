package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s Service) registerLessonRoutes() {
	baseUrl := "/lessons"

	s.Mux.HandleFunc(baseUrl, s.getLessonsHandler).Methods("GET")
	s.Mux.HandleFunc(baseUrl, s.AddLessonHandler).Methods("POST")
}

func (s Service) getLessonsHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("GetLessonsHandler")
	lessons, err := s.Queries.FindLessons(r.Context())
	if err != nil {
		log.Error("Errors occured", err)
	}

	jsonResp, err := json.Marshal(lessons)

	if err != nil {
		log.Error("Errors occured", err)
	}

	w.Write(jsonResp)
}

func (s Service) AddLessonHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("AddLessonHandler")
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		LessonName string `json:"lesson_name"`
	}{}
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Errors occured", err)
	}
	createdLesson, err := s.Queries.InsertLesson(r.Context(), payload.LessonName)
	if err != nil {
		log.Error("Errors occured", err)
	}

	jsonResp, err := json.Marshal(createdLesson)

	if err != nil {
		log.Error("Errors occured", err)
	}

	w.Write(jsonResp)
}
