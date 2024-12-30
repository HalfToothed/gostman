package cmd

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

func load(data Request, model *Model) {

	model.id = data.Id
	model.nameField.SetValue(data.Name)
	model.urlField.SetValue(data.URL)
	model.methodField.SetValue(data.Method)
	model.tabContent[0].SetValue(data.Body)
	model.tabContent[1].SetValue(data.QueryParams)
	model.tabContent[2].SetValue(data.Headers)
	model.response = data.Response
	model.responseViewport.SetContent(model.response)
}

func delete(id string) error {

	// Read the JSON file
	file, err := os.ReadFile(jsonfilePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse the JSON into a slice
	var requests []Request
	if len(file) > 0 {
		if err := json.Unmarshal(file, &requests); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	}

	// Find the item to delete
	index := -1
	for i, req := range requests {
		if req.Id == id {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("id not found: %s", id)
	}

	// Delete the item from the slice
	requests = append(requests[:index], requests[index+1:]...)

	// Write the updated slice back to the file
	updatedData, err := json.MarshalIndent(requests, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	if err := os.WriteFile(jsonfilePath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func checkFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return errors.Is(err, os.ErrNotExist)
}
