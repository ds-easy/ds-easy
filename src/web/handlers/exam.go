package handlers

import (
	"context"
	"ds-easy/src/database/repository"
	utils "ds-easy/src/web/handlers/util"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s Service) RegisterExamRoutes() {
	baseUrl := "/exams"

	s.Mux.HandleFunc(baseUrl, s.getExamsHandler).Methods("GET")
	s.Mux.HandleFunc(baseUrl, s.generateExamHandler).Methods("POST")

	s.Mux.HandleFunc(baseUrl+"/test", s.testHandler).Methods("GET")
}

func (s Service) testHandler(w http.ResponseWriter, r *http.Request) {
	record, err := utils.GetRecordInfo(EXO_FILES, "i2ik67wfjx6l59j")
	if err != nil {
		log.Error("is erroro ", err)
		return
	}

	body, err := utils.DownloadFromPocketBase(EXO_FILES, record.ID, record.File)
	if err != nil {
		log.Error("is erroro ", err)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=file.png")
	w.Write(body)

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

	exoParams := repository.FindRandomExercisesByLessonNameWithLimitParams{
		LessonName: payload.LessonName,
		Limit:      payload.Limit,
	}

	insertParams := repository.InsertExamParams{
		DateOfPassing: payload.DateOfPassing,
		ExamNumber:    payload.ExamNumber,
		ProfessorID:   payload.ProfessorID,
	}

	exam, err := generateExam(s.Queries, exoParams, insertParams, payload.TemplateName)
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	jsonResp, err := json.Marshal(exam)

	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	w.Write(jsonResp)
}

func generateExam(q repository.Queries,
	exoParams repository.FindRandomExercisesByLessonNameWithLimitParams,
	insertExamParams repository.InsertExamParams,
	templateName string) (repository.Exam, error) {
	examExercises, err := q.FindRandomExercisesByLessonNameWithLimit(context.TODO(), exoParams)
	if err != nil {
		log.Error("Errors occured", err)
		return repository.Exam{}, err
	}

	template, err := q.FindTemplateByName(context.TODO(), templateName)
	if err != nil {
		log.Error("Errors occured", err)
		return repository.Exam{}, err
	}

	insertExamParams.TemplateID = template.ID

	exam, err := q.InsertExam(context.TODO(), insertExamParams)
	if err != nil {
		log.Error("Errors occured", err)
		return repository.Exam{}, err
	}

	for _, v := range examExercises {
		err = q.InsertExamExercise(context.TODO(), repository.InsertExamExerciseParams{
			ExamID:     exam.ID,
			ExerciseID: v.ID,
		})
		if err != nil {
			log.Error("Errors occured", err)
			return repository.Exam{}, err
		}
	}

	return exam, nil
}
