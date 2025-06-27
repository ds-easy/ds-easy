package utils

import (
	"bytes"
	"ds-easy/src/database/repository"
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

func PBAddUser(u repository.AddUserParams, password string) (string, error) {
	url := "http://127.0.0.1:8090/api/collections/users/records"
	log.Info("accessing ... ", url)

	user := struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"passwordConfirm"`
		Name            string `json:"name,omitempty"`
	}{
		Email:           u.Email,
		Password:        password,
		PasswordConfirm: password,
		Name:            u.FirstName + " " + u.LastName,
	}

	bodyJson, err := json.Marshal(user)
	if err != nil {
		log.Error("Error marshaling JSON:", err)
		return "", err
	}
	body := bytes.NewBuffer(bodyJson)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Error("error creating POST request: ", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

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
		log.Error("Error reading response body:", err)
		return "", err
	}

	log.Debugln("Response Body: ", string(bodyBytes))

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
		log.Error("Error unmarshaling JSON:", err)
		return "", err
	}

	// Return the "id" from the response
	return response.Id, nil
}

func PBGetUserId(jwt string) (string, error) {
	url := "http://127.0.0.1:8090/api/collections/users/auth-refresh"
	log.Info("accessing ... ", url)
	body := &bytes.Buffer{}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Error("error creating POST request: ", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", jwt)

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
		log.Error("Error reading response body:", err)
		return "", err
	}

	log.Debugln("Response Body: ", string(bodyBytes))

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Error("bad response status: ", resp.Status)
		return "", err
	}

	response := struct {
		Record struct {
			Id string `json:"id"`
		} `json:"record"`
	}{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		log.Error("Error unmarshaling JSON:", err)
		return "", err
	}

	// Return the "id" from the response
	return response.Record.Id, nil
}

func PBUploadFile(file multipart.File, fileName, collection string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		log.Error("error creating form file: ", err)
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Error("error copying file: ", err)
		return "", err
	}

	err = writer.Close()
	if err != nil {
		log.Error("error closing writer: ", err)
		return "", err
	}

	url := os.Getenv("PB_URL") + "collections/" + collection + "/records"

	log.Info("accessing ... ", url)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Error("error creating POST request: ", err)
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("error sending POST request: ", err)
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response body:", err)
		return "", err
	}

	log.Debugln("Response Body:", string(bodyBytes))

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
		log.Error("Error unmarshaling JSON:", err)
		return "", err
	}

	// Return the "id" from the response
	return response.Id, nil
}

func PBCheckPassword(u repository.User, password string) (string, error) {
	log.Info("PBCheckPassword")

	url := "http://127.0.0.1:8090/api/collections/users/auth-with-password"
	log.Info("accessing ... ", url)

	creds := struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}{
		Identity: u.Email,
		Password: password,
	}

	bodyJson, err := json.Marshal(creds)
	if err != nil {
		log.Error("Error marshaling JSON:", err)
		return "", err
	}
	body := bytes.NewBuffer(bodyJson)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Error("error creating POST request: ", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

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
		log.Error("Error reading response body:", err)
		return "", err
	}

	log.Debugln("Response Body: ", string(bodyBytes))

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Error("bad response status: ", resp.Status)
		return "", err
	}

	response := struct {
		Id    string `json:"id"`
		Token string `json:"token"`
	}{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		log.Error("Error unmarshaling JSON:", err)
		return "", err
	}

	// Return the "id" from the response
	return response.Token, nil
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

func DownloadFromPocketBase(collectionName, id string) ([]byte, error) {
	recordInfo, err := GetRecordInfo(collectionName, id)
	if err != nil {
		log.Info("Error:", err)
		return nil, err
	}

	url := os.Getenv("PB_URL") + "files/" + collectionName + "/" + id + "/" + recordInfo.File
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
