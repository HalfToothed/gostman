package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
)

type SavedData struct {
	Variables string    `json:"variables"`
	Requests  []Request `json:"requests"`
}

type Config struct {
	CurrentProject string    `json:"currentProject"`
	Projects       []Project `json:"projects"`
}

type Project struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

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

var appFolder = getAppDataPath()
var configPath = filepath.Join(appFolder, "config.json")
var config *Config
var jsonfilePath string

func init() {
	config = loadConfig()
	// Use configured current project or default
	if config.CurrentProject != "" {
		jsonfilePath = filepath.Join(config.CurrentProject, "gostman.json")
	} else {
		// Fall back to app folder
		jsonfilePath = filepath.Join(appFolder, "gostman.json")
	}
}

func getAppDataPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "Gostman")
	}
	// Linux and macOS path: ~/.local/share/Gostman
	return filepath.Join(os.Getenv("HOME"), ".local", "share", "Gostman")
}

func loadConfig() *Config {
	// Ensure app folder exists
	if err := os.MkdirAll(appFolder, os.ModePerm); err != nil {
		fmt.Println("Failed to create app directory:", err)
	}

	// Check if config file exists
	if checkFileExists(configPath) {
		// Create default config
		cfg := &Config{
			CurrentProject: "",
			Projects:       []Project{},
		}
		saveConfig(cfg)
		return cfg
	}

	// Load existing config
	file, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error reading config:", err)
		// Return default config on error
		return &Config{
			CurrentProject: "",
			Projects:       []Project{},
		}
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		fmt.Println("Error parsing config:", err)
		// Return default config on error
		return &Config{
			CurrentProject: "",
			Projects:       []Project{},
		}
	}

	return &cfg
}

func saveConfig(cfg *Config) {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Println("Error encoding config:", err)
		return
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		fmt.Println("Error saving config:", err)
	}
}

func addProjectIfNotExists(name, path string) {
	for _, project := range config.Projects {
		if project.Path == path {
			return // Already exists
		}
	}
	config.Projects = append(config.Projects, Project{
		Name: name,
		Path: path,
	})
}

func SetCurrentProject(projectPath string) error {
	// Validate the path
	if projectPath == "" {
		return fmt.Errorf("project path cannot be empty")
	}
	
	// Ensure the directory exists
	if err := os.MkdirAll(projectPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	gostmanFile := filepath.Join(projectPath, "gostman.json")
	
	// If file doesn't exist, create it with empty structure
	if checkFileExists(gostmanFile) {
		emptyData := SavedData{
			Variables: "",
			Requests:  []Request{},
		}
		data, err := json.MarshalIndent(emptyData, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to create initial data: %w", err)
		}
		if err := os.WriteFile(gostmanFile, data, 0644); err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
	}
	
	config.CurrentProject = projectPath
	jsonfilePath = gostmanFile
	
	// Add to projects if not already present
	addProjectIfNotExists(filepath.Base(projectPath), projectPath)
	
	saveConfig(config)
	return nil
}

func GetCurrentProject() string {
	return config.CurrentProject
}

func GetProjects() []Project {
	return config.Projects
}

func RemoveProject(projectPath string) error {
	// Don't allow removing the current project if it's the only one
	if len(config.Projects) <= 1 {
		return fmt.Errorf("cannot remove the only project")
	}
	
	// Find and remove the project
	found := false
	newProjects := make([]Project, 0, len(config.Projects)-1)
	for _, p := range config.Projects {
		if p.Path != projectPath {
			newProjects = append(newProjects, p)
		} else {
			found = true
		}
	}
	
	if !found {
		return fmt.Errorf("project not found in configuration")
	}
	
	config.Projects = newProjects
	
	// If we removed the current project, switch to the first available project
	if config.CurrentProject == projectPath {
		if len(config.Projects) > 0 {
			config.CurrentProject = config.Projects[0].Path
			jsonfilePath = filepath.Join(config.CurrentProject, "gostman.json")
		} else {
			config.CurrentProject = ""
			jsonfilePath = filepath.Join(appFolder, "gostman.json")
		}
	}
	
	saveConfig(config)
	return nil
}

func CreateProjectInCurrentDir() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	
	// Use directory name as project name
	projectName := filepath.Base(cwd)
	
	// Ensure the project is added with proper name
	if err := SetCurrentProject(cwd); err != nil {
		return err
	}
	
	// Update the project name in config
	for i, project := range config.Projects {
		if project.Path == cwd {
			config.Projects[i].Name = projectName
			break
		}
	}
	saveConfig(config)
	
	return nil
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

	var saved_data SavedData
	json.Unmarshal(file, &saved_data)

	savedRequests := saved_data.Requests

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

	saved_data.Requests = savedRequests

	updatedData, err := json.MarshalIndent(saved_data, "", " ")

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
	var saved_data SavedData
	if len(file) > 0 {
		if err := json.Unmarshal(file, &saved_data); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	}

	requests := saved_data.Requests

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
	saved_data.Requests = append(requests[:index], requests[index+1:]...)

	// Write the updated slice back to the file
	updatedData, err := json.MarshalIndent(saved_data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	if err := os.WriteFile(jsonfilePath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func SaveVariables(variableString string) string {

	if checkFileExists(jsonfilePath) {

		if err := os.MkdirAll(appFolder, os.ModePerm); err != nil {
			fmt.Println("Failed to create directory:", err)
		}

		myfile, err := os.Create(jsonfilePath)
		if err != nil {
			fmt.Println("Failed to create file:", err)
		}

		defer myfile.Close()
	}

	file, err := os.ReadFile(jsonfilePath)

	if err != nil {
		panic(err)
	}

	var saved_data SavedData
	json.Unmarshal(file, &saved_data)

	var variables map[string]string
	er := json.Unmarshal([]byte(variableString), &variables)
	if er != nil {
		return "Error parsing Environment Variables, JSON structure is incorrect"
	}

	saved_data.Variables = variableString

	updatedData, err := json.MarshalIndent(saved_data, "", " ")

	if err != nil {
		log.Fatalf("Error encoding JSON data: %v", err)
	}

	if err := os.WriteFile(jsonfilePath, updatedData, 0644); err != nil {
		fmt.Println("failed to add the request")
	}

	return "Environment Variables Saved Sucessfully"

}

func checkFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return errors.Is(err, os.ErrNotExist)
}
