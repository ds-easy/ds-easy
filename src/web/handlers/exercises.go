package handlers

import (
	"context"
	"ds-easy/src/database/repository"
	utils "ds-easy/src/web/handlers/util"
	"encoding/json"
	"mime/multipart"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s Service) registerExerciseRoutes() {
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
	file, _, err := r.FormFile("exo_file")
	exerciseName := r.FormValue("exercise_name")
	lessonName := r.FormValue("lesson_name")
	uploadedBy := r.FormValue("uploadedBy")
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}
	defer file.Close()

	insertedExercise, err := s.insertExercise(file, exerciseName, lessonName, uploadedBy)
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	jsonResp, err := json.Marshal(insertedExercise)

	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	w.Write(jsonResp)
}

func (s Service) insertExercise(file multipart.File, exerciseName, lessonName, uploadedBy string) (repository.Exercise, error) {
	lesson, err := s.Queries.FindLessonByName(context.TODO(), lessonName)
	if err != nil {
		log.Error("Errors occured ", err)
		return repository.Exercise{}, err
	}

	user, err := s.Queries.FindUserByEmail(context.TODO(), uploadedBy)
	if err != nil {
		log.Error("Errors occured ", err)
		return repository.Exercise{}, err
	}

	pb_id, err := utils.UploadToPocketBase(file, exerciseName, EXO_FILES)
	if err != nil {
		log.Error("Errors occured ", err)
		return repository.Exercise{}, err
	}

	payload := repository.InsertExerciseParams{
		ExerciseName: exerciseName,
		ExercisePath: pb_id,
		LessonID:     lesson.ID,
		UploadedBy:   user.ID,
	}

	InsertedExercise, err := s.Queries.InsertExercise(context.TODO(), payload)
	if err != nil {
		log.Error("Errors occured ", err)
		return repository.Exercise{}, err
	}
	return InsertedExercise, nil
}
