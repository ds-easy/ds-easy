package handlers

import (
	"ds-easy/src/database/repository"

	"github.com/gorilla/mux"
)

type Service struct {
	Queries repository.Queries
	Mux     *mux.Router
}

// Variables storing collection names for files, basically acting as enums.
var (
	EXO_FILES = "exo_files"
	TEMPLATE  = "template"
)

func (s Service) RegisterRoutes() {
	s.RegisterExamRoutes()
	s.registerExerciseRoutes()
	s.registerLessonRoutes()
	s.registerTemplateRoutes()
	s.registerUserRoutes()
}
