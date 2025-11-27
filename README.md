# GoCycled

GoCycled is a fast, lightweight, and safe trash/recycle bin utility for Linux systems, built in Go. It serves as a modern replacement for the npm `trash-cli` package, providing the same core functionality with improved performance and safety guarantees.

## Features

✨ **Core Functionality**
- Move files to trash instead of permanent deletion
- List all items in trash with details
- Restore files to their original locations
- Empty trash with confirmation
- Remove specific items permanently
- View trash size

🛡️ **Safety First**
- No data loss - files are moved, not deleted
- Unique timestamped filenames prevent conflicts
- Confirmation prompts for destructive operations
- Original path tracking for accurate restoration

⚙️ **Configuration**
- JSON-based config file at `~/.trashrc`
- All settings editable via CLI commands
- Per-user configuration

🚀 **Performance**
- Lightweight and fast
- Minimal dependencies
- Efficient file operations

## Installation

### From Source

```bash
git clone https://github.com/cj3636/GoCycled.git
cd GoCycled
go build -o rc ./cmd/rc
sudo mv rc /usr/local/bin/
```

### Prerequisites

- Go 1.18 or later

## Usage

### Basic Commands

```bash
# Move files to trash
rc put file.txt
rc trash file1.txt file2.txt
rc rm document.pdf

# List items in trash
rc list
rc ls

# Restore files
rc restore                    # Interactive selection
rc restore /path/to/file.txt  # Restore specific file

# Empty trash
rc empty

# Remove specific item permanently
rc remove file.txt

# View trash size
rc size

# Show version
rc version

# Show help
rc help
```

### Configuration Commands

```bash
# View all configuration
rc config

# Get a specific setting
rc config get trash_dir
rc config get confirm_delete

# Set a configuration value
rc config set confirm_delete true
rc config set trash_dir ~/.local/share/Trash

# Reset to default configuration
rc config reset
```

### Configuration Options

The configuration file is stored at `~/.trashrc` in JSON format:

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `trash_dir` | string | `~/.local/share/Trash` | Location of trash directory |
| `confirm_delete` | bool | `true` | Confirm before permanent deletion |
| `auto_empty_days` | int | `30` | Auto-empty trash after N days (future feature) |
| `max_trash_size_mb` | int | `1024` | Maximum trash size in MB (future feature) |

### Examples

```bash
# Replace rm with rc in your workflow
alias rm='rc put'

# Trash multiple files
rc put *.log

# Interactive restore
rc restore

# Empty trash with confirmation
rc empty
```

## Command Structure

The main command is `rc` (recycle), following these patterns:

- `rc <action> [arguments]` - Perform an action
- `rc config <subcommand>` - Manage configuration

Supported actions:
- `put`, `trash`, `rm` - Move to trash
- `list`, `ls` - List items
- `restore` - Restore items
- `empty` - Empty trash
- `remove`, `delete` - Permanently delete item
- `size` - Show trash size

## Architecture

```
GoCycled/
├── cmd/rc/           # Main CLI application
│   └── main.go       # Command-line interface
├── pkg/
│   ├── config/       # Configuration management
│   ├── trash/        # Trash operations
│   └── ui/           # User interface
└── go.mod            # Go module file
```

### Trash Storage

Files are stored in `~/.local/share/Trash/` (XDG compliant):
- `files/` - Actual trashed files with timestamped names
- `info/` - JSON metadata files tracking original paths and deletion times

## Safety Features

1. **No Data Loss**: Files are moved, not deleted immediately
2. **Unique Names**: Timestamped filenames prevent overwrites
3. **Metadata Tracking**: JSON files track original paths and timestamps
4. **Confirmation Prompts**: Optional confirmations for destructive operations
5. **Original Path Restoration**: Files can be restored to exact original locations

## Comparison with trash-cli

| Feature | GoCycled | trash-cli |
|---------|----------|-----------|
| Language | Go | Node.js |
| Performance | Fast | Slower |
| Dependencies | Minimal | npm ecosystem |
| Configuration | `~/.trashrc` | Various |
| Command | `rc` | `trash`, `trash-list`, etc. |
| UI Options | Basic | Basic |
| Config Management | Built-in CLI | Manual editing |

## Development

### Building

```bash
go build -o rc ./cmd/rc
```

### Testing

```bash
go test -v ./...
```

### Running Tests with Coverage

```bash
go test -cover ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- Inspired by [trash-cli](https://github.com/sindresorhus/trash-cli)
