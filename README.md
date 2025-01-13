# 🖥️ Gostman

A terminal-based API client built with Bubble Tea for creating, sending, and managing HTTP requests in an interactive and user-friendly way.

![gostman](https://github.com/user-attachments/assets/98ad0be6-479e-432f-8794-bd495a401872)

## ✨ Features

- Create and send HTTP requests (GET, POST, PUT, DELETE, etc.).
- Save, load, and manage requests.
- Edit and delete saved requests easily.
- Dynamic UI with support for status messages, and detailed responses.

## 📥 Installation 

1. Ensure you have [Go](https://go.dev/) installed.
2. Clone this repository:
   ```bash
   git clone https://github.com/HalfToothed/gostman.git
   cd gostman
   ```

3. Run the project:
    ```bash
    go run .
    ```

## 🧑‍💻 Usage 

### Use keyboard shortcuts to interact with the UI:

- Ctrl+C: Quit the application.
- Shift+Tab: Move Around
- Shift+Up: Send a request.
- Ctrl+S: Save the current request.
- Shift+Left: Load a saved request.

### Saving and Loading Requests 

Requests are saved as JSON files in the user's home directory under a dedicated folder. The JSON file structure allows for efficient updates and deletions.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues or pull requests to improve this project. 🙌
