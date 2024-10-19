package handlers

import (
	"bytes"
	"context"
	"ds-easy/src/database/repository"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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
	pb_id, err := uploadToPocketBase(file, exerciseName)
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
	}

	InsertedExercise, err := s.Queries.InsertExercise(context.TODO(), payload)
	if err != nil {
		log.Error("Errors occured ", err)
		return repository.Exercise{}, err
	}
	return InsertedExercise, nil
}

func uploadToPocketBase(file multipart.File, exerciseName string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("files", exerciseName)
	if err != nil {
		log.Error("error creating form file: ", err)
	}

	// Copy the file contents to the form field
	_, err = io.Copy(part, file)
	if err != nil {
		log.Error("error copying file: ", err)
	}

	// Close the multipart writer to finalize the request
	err = writer.Close()
	if err != nil {
		log.Error("error closing writer: ", err)
	}

	// Create the POST request with the multipart form data
	req, err := http.NewRequest("POST", "http://127.0.0.1:8090/api/collections/exo_files/records", body)
	if err != nil {
		log.Error("error creating POST request: ", err)
	}

	// Set the content type for the multipart form data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("error sending POST request: ", err)
	}
	defer resp.Body.Close()

	log.Info(resp.StatusCode)

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Error("bad response status: ", resp.Status)
	}

	fmt.Println("File uploaded successfully with status:", resp.Status)

	response := struct {
		Id string `json:"id"`
	}{}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return "", err
	}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		log.Println("Error unmarshaling JSON:", err)
		return "", err
	}

	// Return the "id" from the response
	return response.Id, nil
}
