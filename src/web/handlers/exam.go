package handlers

import (
	"ds-easy/src/database/repository"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s Service) RegisterExamRoutes() {
	baseUrl := "/exams"

	s.Mux.HandleFunc(baseUrl, s.getExamsHandler).Methods("GET")
	s.Mux.HandleFunc(baseUrl, s.generateExamHandler).Methods("POST")
}

func (s Service) getExamsHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("getExamsHandler")
	exams, err := s.Queries.FindExams(r.Context())
	if err != nil {
		log.Error("Errors occured querying exams ", err)
		w.WriteHeader(500)
		return
	}

	jsonResp, err := json.Marshal(exams)

	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	w.Write(jsonResp)
}

func (s Service) generateExamHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("generateExamHandler")

	decoder := json.NewDecoder(r.Body)
	payload := struct {
		repository.FindRandomExercisesByLessonNameWithLimitParams
		repository.InsertExamParams
		TemplateName string `json:"template_name"`
	}{}
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	examExercises, err := s.Queries.FindRandomExercisesByLessonNameWithLimit(r.Context(), repository.FindRandomExercisesByLessonNameWithLimitParams{
		LessonName: payload.LessonName,
		Limit:      payload.Limit,
	})
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	template, err := s.Queries.FindTemplateByName(r.Context(), payload.TemplateName)
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	exam, err := s.Queries.InsertExam(r.Context(), repository.InsertExamParams{
		DateOfPassing: payload.DateOfPassing,
		ExamNumber:    payload.ExamNumber,
		ProfessorID:   payload.ProfessorID,
		TemplateID:    template.ID,
	})
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	for _, v := range examExercises {
		err = s.Queries.InsertExamExercise(r.Context(), repository.InsertExamExerciseParams{
			ExamID:     exam.ID,
			ExerciseID: v.ID,
		})
		if err != nil {
			log.Error("Errors occured", err)
			w.WriteHeader(500)
			return
		}
	}

	jsonResp, err := json.Marshal(exam)

	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	w.Write(jsonResp)
}
