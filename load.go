package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func load(data Request, model *Model) {

	model.id = data.Id
	model.nameField.SetValue(data.Name)
	model.urlField.SetValue(data.URL)
	model.methodField.SetValue(data.Method)
	model.tabContent[0].SetValue(data.Body)
	model.tabContent[1].SetValue(data.QueryParams)
	model.tabContent[2].SetValue(data.Headers)
	model.response = data.Response
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
