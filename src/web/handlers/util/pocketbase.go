package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

type ResponseBody struct {
	CollectionID   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	Created        string `json:"created"`
	File           string `json:"file"`
	ID             string `json:"id"`
	Updated        string `json:"updated"`
}

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
	url := os.Getenv("PB_URL") + "collections/" + collection + "/records"

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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return "", err
	}

	log.Println("Response Body:", string(bodyBytes))

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Error("bad response status: ", resp.Status)
		return "", err
	}

	response := struct {
		Id string `json:"id"`
	}{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		log.Println("Error unmarshaling JSON:", err)
		return "", err
	}

	// Return the "id" from the response
	return response.Id, nil
}

func GetRecordInfo(collectionName, id string) (ResponseBody, error) {
	url := os.Getenv("PB_URL") + "collections/" + collectionName + "/records/" + id

	resp, err := http.Get(url)
	if err != nil {
		return ResponseBody{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseBody{}, err
	}

	var responseBody ResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return ResponseBody{}, err
	}

	return responseBody, nil
}

func DownloadFromPocketBase(collectionName, id, fileName string) ([]byte, error) {
	url := os.Getenv("PB_URL") + "files/" + collectionName + "/" + id + "/" + fileName
	log.Info("accessing ... ", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Info("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Info("Error:", err)
		return nil, err
	}

	return body, nil
}
