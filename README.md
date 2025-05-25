# 🖥️ Gostman

A terminal-based API client built with Bubble Tea for creating, sending, and managing HTTP requests in an interactive and user-friendly way.

![gostman](https://github.com/user-attachments/assets/65c46e9d-2600-47c9-809f-779b5531f023)

## ✨ Features

- Create and send HTTP requests (GET, POST, PUT, DELETE, etc.).
- Save, load, and manage requests.
- Edit and delete saved requests easily.
- Dynamic UI with support for status messages, and detailed responses.

## 📥 Install

### Go

_If you have Go already, install the executable yourself_

1. Run the following command:
   ```bash
   go install github.com/halftoothed/gostman@latest
   ```
2. The tool is ready to use!
   ```bash
   gostman
   ```

### Nix

If you are a Nixos user, you can try gostman via flake:

```nix
nix run github:HalfToothed/gostman
```

Or add it to your system:

#### flake.nix

```nix
{
  inputs = {
    ...
    # gostman input
    gostman.url = "github:HalfToothed/gostman";
    # optionally follow your nixpkgs version
    gostman.inputs.nixpkgs.follows = "nixpkgs";
  };
}
```

#### configuration.nix

```nix
{ inputs, pkgs, ...}: {
  environment.systemPackages =
    let gostman = inputs.gostman.packages.${pkgs.system}.default;
    in [ gostman ];
}
```

Now you can rebuild your system and run gostman

```bash
sudo nixos-rebuild switch && gostman
```

## 🧑‍💻 Usage

### Use keyboard shortcuts to interact with the UI:

- Ctrl + C: Quit the application.
- Tab: Move Around
- Ctrl + Arrow Keys: Change Tabs (Body/Param/Header)
- Enter: Send a request.
- Ctrl + S: Save the current request.
- Ctrl + E: Open Environment Variables page
- Ctrl + D: Open Dashboard.
- Ctrl + H: Open Help page

### Saving and Loading Requests

Requests are saved as JSON files in the user's home directory under a dedicated folder. The JSON file structure allows for efficient updates and deletions.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues or pull requests to improve this project. 🙌
