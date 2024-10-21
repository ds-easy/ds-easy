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

func (s Service) RegisterTemplateRoutes() {
	baseUrl := "/templates"

	s.Mux.HandleFunc(baseUrl, s.getTemplatesHandler).Methods("GET")
	s.Mux.HandleFunc(baseUrl, s.InsertTemplateHandler).Methods("POST")
}

func (s Service) getTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	templates, err := s.Queries.FindTemplates(r.Context())
	if err != nil {
		log.Error("Errors occured", err)
	}

	jsonResp, err := json.Marshal(templates)

	if err != nil {
		log.Error("Errors occured", err)
	}

	w.Write(jsonResp)
}

func (s Service) InsertTemplateHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("InsertTemplateHandler")
	file, _, err := r.FormFile("template_file")
	uploadedBy := r.FormValue("uploadedBy")
	templateName := r.FormValue("template_name")
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}
	defer file.Close()

	insertedTemplate, err := insertTemplate(s.Queries, file, templateName, uploadedBy)
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	jsonResp, err := json.Marshal(insertedTemplate)
	if err != nil {
		log.Error("Errors occured", err)
		w.WriteHeader(500)
		return
	}

	w.Write(jsonResp)
}

func insertTemplate(q repository.Queries, file multipart.File, templateName, uploadedBy string) (repository.Template, error) {
	pb_id, err := utils.UploadToPocketBase(file, templateName, TEMPLATE)
	if err != nil {
		log.Error("Errors occured ", err)
		return repository.Template{}, err
	}

	user, err := q.FindUserByEmail(context.TODO(), uploadedBy)
	if err != nil {
		log.Error("Errors occured ", err)
		return repository.Template{}, err
	}

	payload := repository.InsertTemplateParams{
		PbFileID:     pb_id,
		UploadedBy:   user.ID,
		TemplateName: templateName,
	}

	insertedTemplate, err := q.InsertTemplate(context.TODO(), payload)
	if err != nil {
		log.Error("Errors occured ", err)
		return repository.Template{}, err
	}

	return insertedTemplate, nil
}
