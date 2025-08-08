package handlers

import (
	"encoding/json"
	"net/http"
)

func (s Service) registerHealthRoutes() {
	baseUrl := "/health"

	s.Mux.HandleFunc(baseUrl, s.getHealthStatus).Methods("GET")
}

func (s Service) getHealthStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status": "UP",
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error serializing response", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}
