# 🚀 Route Keeper

A beautiful TUI (Terminal User Interface) application for monitoring and managing API endpoints with real-time status updates and historical data.

![Route Keeper Screenshot](https://github.com/lutefd/route-keeper/raw/main/screenshot.png)

## ✨ Features

- Monitor multiple API endpoints with custom intervals
- View detailed request/response information
- Configure custom headers and query parameters
- Beautiful and intuitive terminal interface
- Light and dark mode support
- Real-time status updates
- Response time tracking

## 📦 Installation

### Using Go install (requires Go 1.16+)

```bash
go install github.com/lutefd/route-keeper/cmd/route-keeper@latest
```

### Download pre-built binaries

1. Visit the [Releases](https://github.com/lutefd/route-keeper/releases) page
2. Download the appropriate binary for your system
3. Make the binary executable (`chmod +x route-keeper` on Unix-like systems)
4. Move the binary to a directory in your `PATH`

## 🚀 Usage

```bash
# Start the application
route-keeper

# Show version information
route-keeper --version
```

### Keybindings

- **↑/↓/j/k**: Navigate menus and lists
- **Enter**: Select item or confirm action
- **Esc**: Go back or cancel
- **q**: Quit the application
- **s**: Start/stop monitoring
- **e**: Edit profile
- **d**: Delete profile
- **c**: Create new profile

## 🛠 Building from Source

### Prerequisites

- Go 1.16 or later
- Git
- Make (optional, but recommended)

### Using Make (recommended)

1. Clone the repository:

   ```bash
   git clone https://github.com/lutefd/route-keeper.git
   cd route-keeper
   ```

2. Build the application:

   ```bash
   # Build the application
   make build

   # The binary will be available at bin/route-keeper
   ```

3. (Optional) Install to your GOPATH:
   ```bash
   make install
   ```

### Manual Build

If you don't have Make installed, you can build manually:

```bash
git clone https://github.com/lutefd/route-keeper.git
cd route-keeper

go build -o route-keeper ./cmd/route-keeper
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 Version Information

Version information is automatically populated during the build process. To view the current version information:

```bash
# For installed version
route-keeper --version

# For local build
make build && ./bin/route-keeper --version
```

The version information includes:

- Application version (from git tag)
- Git commit hash
- Build timestamp
- Go version used for building
- Builder's username

### Releasing New Versions

1. Create a new git tag:

   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. Build release binaries:

   ```bash
   make release
   ```

   This will create platform-specific binaries in the `dist/` directory.

3. Create a new GitHub release and upload the binaries.
