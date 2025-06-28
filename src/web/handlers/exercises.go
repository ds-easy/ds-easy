package handlers

import (
	"context"
	"database/sql"
	"ds-easy/src/database/repository"
	utils "ds-easy/src/web/handlers/util"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (s Service) registerExerciseRoutes() {
	baseUrl := "/exercises"

	s.Mux.HandleFunc(baseUrl, s.getExercisesHandler).Methods("GET")
	s.Mux.HandleFunc(baseUrl, s.addExerciseHandler).Methods("POST")
	s.Mux.HandleFunc(baseUrl+"/public", s.getPublicExercisesHandler).Methods("GET")
	s.Mux.HandleFunc(baseUrl+"/accessible", s.getAccessibleExercisesHandler).Methods("GET")
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

func (s Service) getPublicExercisesHandler(w http.ResponseWriter, r *http.Request) {
	exercises, err := s.Queries.FindPublicExercises(r.Context())
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	jsonResp, err := json.Marshal(exercises)
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	w.Write(jsonResp)
}

func (s Service) getAccessibleExercisesHandler(w http.ResponseWriter, r *http.Request) {
	userID := s.getUserIDFromQuery(r).Int64

	exercises, err := s.Queries.FindAccessibleExercises(r.Context(), userID)
	if err != nil {
		log.Error("Error finding accessible exercises: ", err)
		http.Error(w, "Failed to fetch exercises", http.StatusInternalServerError)
		return
	}

	s.writeJSONResponse(w, exercises)
}

func (s Service) addExerciseHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("exo_file")
	exerciseName := r.FormValue("exercise_name")
	lessonName := r.FormValue("lesson_name")
	uploadedBy := r.FormValue("uploadedBy")
	isPublicStr := r.FormValue("is_public")

	isPublic := false
	if isPublicStr != "" {
		if parsed, err := strconv.ParseBool(isPublicStr); err == nil {
			isPublic = parsed
		}
	}

	if err != nil {
		log.Error("Errors occured ", err)
		w.WriteHeader(500)
		return
	}
	defer file.Close()

	insertedExercise, err := s.insertExercise(file, exerciseName, lessonName, uploadedBy, isPublic)
	if err != nil {
		log.Error("Errors occured ", err)
		w.WriteHeader(500)
		return
	}

	jsonResp, err := json.Marshal(insertedExercise)

	if err != nil {
		log.Error("Errors occured ", err)
		w.WriteHeader(500)
		return
	}

	w.Write(jsonResp)
}

func (s Service) insertExercise(file multipart.File, exerciseName, lessonName, uploadedBy string, isPublic bool) (repository.Exercise, error) {
	pb_id, err := utils.PBUploadFile(file, exerciseName, EXO_FILES)
	if err != nil {
		log.Error("Errors occured ", err)
		return repository.Exercise{}, err
	}

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

	payload := repository.InsertExerciseParams{
		ExerciseName: exerciseName,
		ExercisePath: pb_id,
		LessonID:     lesson.ID,
		UploadedBy:   user.ID,
		IsPublic:     isPublic,
	}

	InsertedExercise, err := s.Queries.InsertExercise(context.TODO(), payload)
	if err != nil {
		log.Error("Errors occured ", err)
		return repository.Exercise{}, err
	}
	return InsertedExercise, nil
}

// /////////// HELPER FUNCTIONS //////////////
func (s Service) getUserIDFromQuery(r *http.Request) sql.NullInt64 {
	userIDStr := r.URL.Query().Get("uploaded_by")
	if userIDStr == "" {
		return sql.NullInt64{Valid: false}
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return sql.NullInt64{Valid: false}
	}

	return sql.NullInt64{Int64: userID, Valid: true}
}

func (s Service) writeJSONResponse(w http.ResponseWriter, data interface{}) {
	jsonResp, err := json.Marshal(data)
	if err != nil {
		log.Error("Error marshaling JSON: ", err)
		http.Error(w, "Failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}
