package handlers

import (
	"ds-easy/src/database/repository"

	"github.com/gorilla/mux"
)

type Service struct {
	Queries repository.Queries
	Mux     *mux.Router
}
