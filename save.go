package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var appDataPath = os.Getenv("APPDATA")
var appFolder = filepath.Join(appDataPath, "Gostman")
var jsonfilePath = filepath.Join(appFolder, "gostman.json")

// Request represents the structure of a single saved request
type Request struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Method      string `json:"method"`
	Headers     string `json:"headers"`
	Body        string `json:"body"`
	QueryParams string `json:"queryParams"`
	Response    string `json:"response"`
}

func SaveRequests(request Request) {

	if checkFileExists(jsonfilePath) {

		if err := os.MkdirAll(appFolder, os.ModePerm); err != nil {
			fmt.Println("Failed to create directory:", err)
			return
		}

		myfile, err := os.Create(jsonfilePath)
		if err != nil {
			fmt.Println("Failed to create file:", err)
			return
		}

		defer myfile.Close()
	}

	file, err := os.ReadFile(jsonfilePath)

	if err != nil {
		panic(err)
	}

	var savedRequests []Request
	json.Unmarshal(file, &savedRequests)

	if request.Id == "" {

		request.Id = uuid.New().String()
		savedRequests = append(savedRequests, request)

	} else {

		for i, r := range savedRequests {
			if r.Id == request.Id {
				savedRequests[i] = request
				break
			}
		}
	}

	updatedData, err := json.MarshalIndent(savedRequests, "", " ")

	if err != nil {
		log.Fatalf("Error encoding JSON data: %v", err)
	}

	if err := os.WriteFile(jsonfilePath, updatedData, 0644); err != nil {
		fmt.Println("failed to add the request")
	}

}

func save(m Model) {
	// Example request data
	request := Request{
		Id:          m.id,
		Name:        m.nameField.Value(),
		URL:         m.urlField.Value(),
		Method:      m.methodField.Value(),
		Body:        m.tabContent[0].Value(),
		QueryParams: m.tabContent[1].Value(),
		Headers:     m.tabContent[2].Value(),
		Response:    m.response,
	}
	SaveRequests(request)
}

func checkFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return errors.Is(err, os.ErrNotExist)
}
