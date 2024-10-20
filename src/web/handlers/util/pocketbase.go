package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func UploadToPocketBase(file multipart.File, fileName, collection string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		log.Error("error creating form file: ", err)
		return "", err
	}

	// Copy the file contents to the form field
	_, err = io.Copy(part, file)
	if err != nil {
		log.Error("error copying file: ", err)
		return "", err
	}

	// Close the multipart writer to finalize the request
	err = writer.Close()
	if err != nil {
		log.Error("error closing writer: ", err)
		return "", err
	}

	// Create the POST request with the multipart form data
	url := "http://127.0.0.1:8090/api/collections/" + collection + "/records"

	log.Info("accessing ... ", url)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Error("error creating POST request: ", err)
		return "", err
	}

	// Set the content type for the multipart form data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("error sending POST request: ", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Error("bad response status: ", resp.Status)
		return "", err
	}

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
