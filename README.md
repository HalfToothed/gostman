# 🖥️ Gostman

A terminal-based API client built with Bubble Tea for creating, sending, and managing HTTP requests in an interactive and user-friendly way.

![gostman](https://github.com/user-attachments/assets/65c46e9d-2600-47c9-809f-779b5531f023)

## ✨ Features

- Create and send HTTP requests (GET, POST, PUT, DELETE, etc.)
- Save, load, and manage requests
- Edit and delete saved requests easily
- Dynamic UI with support for status messages, and detailed responses
- **Project Management**: Automatically detects and manages API collections per directory
- **Multi-project support**: Track and switch between different projects
- **Environment variables**: Support for dynamic request configuration

## 📥 Install

_If you have Go already, install the executable yourself_

1. Run the following command:
   ```bash
   go install github.com/halftoothed/gostman@latest
   ```
2. The tool is ready to use!
    ```bash
   gostman
   ```

## 🧑‍💻 Usage 

### Project Management

Gostman automatically manages API collections per project directory:
1. When you run `gostman`, it checks for a `gostman.json` file in your current directory
2. If found, it automatically loads that project's requests
3. If not found, it prompts you to create a new project or use an existing one
4. Each directory can have its own collection of API requests

### Keyboard Shortcuts

#### Main Interface
- **Ctrl + C**: Quit the application
- **Tab**: Navigate between fields (Name, Method, URL, Content tabs)
- **Shift + ←/→**: Switch between content tabs (Body/Params/Headers)
- **Enter**: Send the current request
- **Ctrl + S**: Save the current request

#### Project Management
- **Ctrl + P**: Quick project switcher
- **Ctrl + D**: Open main dashboard

#### Dashboard Navigation
- **p**: Toggle project view
- **r**: Remove selected project (when in project view)
- **n**: Create new request
- **d**: Delete selected request
- **Enter**: Load selected request or switch to selected project

#### Other Features
- **Ctrl + E**: Open Environment Variables editor
- **Ctrl + H**: Open Help page
- **Esc**: Go back/cancel current action

### Project Workflow

#### Creating a New Project
1. Navigate to your project directory
2. Run `gostman`
3. Choose "y" to create a new `gostman.json` file
4. Start creating and saving requests

#### Switching Between Projects
1. Use **Ctrl + P** for quick project switching
2. Or use **Ctrl + D** → **p** to browse all projects
3. Select any project to switch to it instantly

#### Managing Requests
- Requests are automatically saved to the current project's `gostman.json`
- Each project maintains its own collection of requests
- Use the dashboard (**Ctrl + D**) to browse, edit, and delete saved requests

### Environment Variables

Use environment variables to make your requests dynamic:
1. Press **Ctrl + E** to open the environment editor
2. Define variables in JSON format: `{"baseUrl": "https://api.example.com"}`
3. Use variables in requests with double braces: `{{baseUrl}}/users`

### Saving and Loading Requests 

Requests are saved as JSON files in your project directory or user's home directory. The JSON file structure allows for efficient updates and deletions. Each project maintains its own `gostman.json` file with requests and environment variables.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues or pull requests to improve this project. 🙌
