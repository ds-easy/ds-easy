package handlers

import (
	"context"
	"ds-easy/src/database/repository"
	utils "ds-easy/src/web/handlers/util"
	"encoding/json"
	"net/http"
	"strings"
	"fmt"

	gotypst "github.com/francescoalemanno/gotypst"
	log "github.com/sirupsen/logrus"
)

type LessonRequest struct {
	LessonName string `json:"lesson_name"`
	Limit      *int64 `json:"limit,omitempty"`
}

type GenerateExamPayload struct {
	Lessons      []LessonRequest `json:"lessons"`
	GlobalLimit  *int64          `json:"global_limit,omitempty"`
	TemplateName string          `json:"template_name"`
	repository.InsertExamParams
}

func (s Service) RegisterExamRoutes() {
	baseUrl := "/exams"

	s.Mux.HandleFunc(baseUrl, s.getExamsHandler).Methods("GET")
	s.Mux.HandleFunc(baseUrl, s.generateExamHandler).Methods("POST")

}

func (s Service) getExamsHandler(w http.ResponseWriter, r *http.Request) {
	exams, err := s.Queries.FindExams(r.Context())
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
	}

	jsonResp, err := json.Marshal(exams)
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
	}

	w.Write(jsonResp)
}

func (s Service) generateExamHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("generateExamHandler")

	decoder := json.NewDecoder(r.Body)
	var payload GenerateExamPayload
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	processedLessons, err := processLessonLimits(payload.Lessons, payload.GlobalLimit)
	if err != nil {
		log.Error("Error processing lesson limits: ", err)
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	insertParams := repository.InsertExamParams{
		DateOfPassing: payload.DateOfPassing,
		ExamNumber:    payload.ExamNumber,
		ProfessorID:   payload.ProfessorID,
	}

	exam, err := generateExamFromMultipleLessons(s.Queries, processedLessons, insertParams, payload.TemplateName)
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Write(exam)
}

func processLessonLimits(lessons []LessonRequest, globalLimit *int64) ([]LessonRequest, error) {
	if globalLimit == nil {
		for _, lesson := range lessons {
			if lesson.Limit == nil {
				return nil, fmt.Errorf("lesson limit is required when global_limit is not provided for lesson: %s", lesson.LessonName)
			}
		}
		return lessons, nil
	}

	if *globalLimit < int64(len(lessons)) {
		return nil, fmt.Errorf("global_limit (%d) cannot be less than the number of lessons (%d). Each lesson needs at least 1 exercise", *globalLimit, len(lessons))
	}
	return distributeGlobalLimit(lessons, *globalLimit), nil
}


func distributeGlobalLimit(lessons []LessonRequest, globalLimit int64) []LessonRequest {
	if len(lessons) == 0 {
		return lessons
	}

	numLessons := int64(len(lessons))
	baseLimit := globalLimit / numLessons
	remainder := globalLimit % numLessons

	processedLessons := make([]LessonRequest, len(lessons))

	for i, lesson := range lessons {
		limit := baseLimit
		if int64(i) < remainder {
			limit++
		}

		processedLessons[i] = LessonRequest{
			LessonName: lesson.LessonName,
			Limit:      &limit,
		}

		log.Infof("Lesson '%s' assigned limit: %d", lesson.LessonName, limit)
	}

	return processedLessons
}

func generateExamFromMultipleLessons(q repository.Queries,
	lessons []LessonRequest,
	insertExamParams repository.InsertExamParams,
	templateName string) ([]byte, error) {

	var allExercises []repository.Exercise
	var lessonNames []string

	for _, lesson := range lessons {
		exoParams := repository.FindRandomAccessibleExercisesByLessonNameWithLimitParams{
			LessonName: lesson.LessonName,
			UploadedBy: insertExamParams.ProfessorID,
			Limit:      *lesson.Limit,
		}

		examExercises, err := q.FindRandomAccessibleExercisesByLessonNameWithLimit(context.TODO(), exoParams)
		if err != nil {
			log.Error("Errors occured", err)
			return nil, err
		}

		allExercises = append(allExercises, examExercises...)
		lessonNames = append(lessonNames, lesson.LessonName)
	}

	template, err := q.FindTemplateByName(context.TODO(), templateName)
	if err != nil {
		log.Error("Errors occured", err)
		return nil, err
	}

	insertExamParams.TemplateID = template.ID

	exam, err := q.InsertExam(context.TODO(), insertExamParams)
	if err != nil {
		log.Error("Errors occured ", err)
		return nil, err
	}

	professor, err := q.FindUserById(context.TODO(), insertExamParams.ProfessorID)
	if err != nil {
		log.Error("Errors occured ", err)
		return nil, err
	}

	sb := strings.Builder{}
	for _, v := range allExercises {
		err = q.InsertExamExercise(context.TODO(), repository.InsertExamExerciseParams{
			ExamID:     exam.ID,
			ExerciseID: v.ID,
		})
		if err != nil {
			log.Error("Errors occured", err)
			return nil, err
		}

		exoFile, err := utils.DownloadFromPocketBase(EXO_FILES, v.ExercisePath)
		if err != nil {
			log.Error("Errors occured", err)
			return nil, err
		}

		_, _ = sb.WriteString("#question[\n\n")
		_, _ = sb.Write(exoFile)
		_, _ = sb.WriteString("\n\n]")
		_, _ = sb.WriteString("\n\n")
	}

	templateFile, err := utils.DownloadFromPocketBase(TEMPLATE, template.PbFileID)
	if err != nil {
		log.Error("Errors occured", err)
		return nil, err
	}

	templateString := string(templateFile)
	allLessonsName := strings.Join(lessonNames, ", ")
	templateString = replaceInfo(templateString, professor, allLessonsName, insertExamParams)

	result := strings.Replace(templateString, "{{EXERCISES}}", sb.String(), 1)

	log.Info("RESULT", result)

	resultPdf, err := gotypst.PDF([]byte(result))
	if err != nil {
		log.Error("Errors occured", err)
		return nil, err
	}

	return resultPdf, nil
}

func replaceInfo(template string, professor repository.User, lessonName string, insertExamParams repository.InsertExamParams) string {
	result := strings.Replace(template, "{{lesson}}", lessonName, 1)
	result = strings.Replace(result, "{{course}}", "", 1)
	result = strings.Replace(result, "{{date}}", insertExamParams.DateOfPassing.Format("02/01/2006"), 1)
	result = strings.Replace(result, "{{duration}}", "2 heures", 1)
	result = strings.Replace(result, "{{prof_name}}", professor.FirstName+" "+professor.LastName, 1)
	result = strings.Replace(result, "{{school_name}}", "", 1)

	return result
}
