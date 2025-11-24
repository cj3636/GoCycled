package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cj3636/GoCycled/pkg/config"
	"github.com/cj3636/GoCycled/pkg/trash"
	"github.com/cj3636/GoCycled/pkg/ui"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Create UI
	var userUI ui.UI
	if cfg.UseFancyUI {
		userUI = ui.NewFancyUI()
	} else {
		userUI = ui.NewBasicUI()
	}

	// Create trash manager
	trashMgr, err := trash.NewManager(cfg.TrashDir)
	if err != nil {
		userUI.Error(fmt.Sprintf("Failed to initialize trash manager: %v", err))
		os.Exit(1)
	}

	// Parse command
	command := os.Args[1]

	switch command {
	case "put", "trash", "rm":
		cmdPut(trashMgr, userUI, cfg, os.Args[2:])
	case "list", "ls":
		cmdList(trashMgr, userUI)
	case "restore":
		cmdRestore(trashMgr, userUI, os.Args[2:])
	case "empty":
		cmdEmpty(trashMgr, userUI, cfg)
	case "remove", "delete":
		cmdRemove(trashMgr, userUI, os.Args[2:])
	case "size":
		cmdSize(trashMgr, userUI)
	case "config":
		cmdConfig(cfg, userUI, os.Args[2:])
	case "version", "--version", "-v":
		fmt.Printf("rc version %s\n", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		userUI.Error(fmt.Sprintf("Unknown command: %s", command))
		printUsage()
		os.Exit(1)
	}
}

func cmdPut(trashMgr *trash.Manager, userUI ui.UI, cfg *config.Config, args []string) {
	if len(args) == 0 {
		userUI.Error("No files specified")
		os.Exit(1)
	}

	for _, path := range args {
		if err := trashMgr.Put(path); err != nil {
			userUI.Error(fmt.Sprintf("Failed to trash %s: %v", path, err))
		} else {
			userUI.Success(fmt.Sprintf("Moved to trash: %s", path))
		}
	}
}

func cmdList(trashMgr *trash.Manager, userUI ui.UI) {
	items, err := trashMgr.List()
	if err != nil {
		userUI.Error(fmt.Sprintf("Failed to list trash: %v", err))
		os.Exit(1)
	}

	userUI.DisplayItems(items)
}

func cmdRestore(trashMgr *trash.Manager, userUI ui.UI, args []string) {
	items, err := trashMgr.List()
	if err != nil {
		userUI.Error(fmt.Sprintf("Failed to list trash: %v", err))
		os.Exit(1)
	}

	if len(items) == 0 {
		userUI.Info("Trash is empty")
		return
	}

	var trashName string

	if len(args) > 0 {
		// Restore by original path or trash name
		targetPath := args[0]
		for _, item := range items {
			if item.OriginalPath == targetPath || filepath.Base(item.TrashPath) == targetPath {
				trashName = filepath.Base(item.TrashPath)
				break
			}
		}
		if trashName == "" {
			userUI.Error(fmt.Sprintf("Item not found: %s", targetPath))
			os.Exit(1)
		}
	} else {
		// Interactive selection
		selected, err := userUI.SelectItem(items)
		if err != nil {
			userUI.Error(fmt.Sprintf("Selection failed: %v", err))
			os.Exit(1)
		}
		trashName = filepath.Base(selected)
	}

	if err := trashMgr.Restore(trashName); err != nil {
		userUI.Error(fmt.Sprintf("Failed to restore: %v", err))
		os.Exit(1)
	}

	userUI.Success("Item restored successfully")
}

func cmdEmpty(trashMgr *trash.Manager, userUI ui.UI, cfg *config.Config) {
	items, err := trashMgr.List()
	if err != nil {
		userUI.Error(fmt.Sprintf("Failed to list trash: %v", err))
		os.Exit(1)
	}

	if len(items) == 0 {
		userUI.Info("Trash is already empty")
		return
	}

	if cfg.ConfirmDelete {
		if !userUI.Confirm(fmt.Sprintf("Permanently delete %d items?", len(items))) {
			userUI.Info("Operation cancelled")
			return
		}
	}

	if err := trashMgr.Empty(); err != nil {
		userUI.Error(fmt.Sprintf("Failed to empty trash: %v", err))
		os.Exit(1)
	}

	userUI.Success(fmt.Sprintf("Permanently deleted %d items", len(items)))
}

func cmdRemove(trashMgr *trash.Manager, userUI ui.UI, args []string) {
	if len(args) == 0 {
		userUI.Error("No item specified")
		os.Exit(1)
	}

	items, err := trashMgr.List()
	if err != nil {
		userUI.Error(fmt.Sprintf("Failed to list trash: %v", err))
		os.Exit(1)
	}

	targetPath := args[0]
	var trashName string

	for _, item := range items {
		if item.OriginalPath == targetPath || filepath.Base(item.TrashPath) == targetPath {
			trashName = filepath.Base(item.TrashPath)
			break
		}
	}

	if trashName == "" {
		userUI.Error(fmt.Sprintf("Item not found: %s", targetPath))
		os.Exit(1)
	}

	if err := trashMgr.Remove(trashName); err != nil {
		userUI.Error(fmt.Sprintf("Failed to remove: %v", err))
		os.Exit(1)
	}

	userUI.Success("Item permanently deleted")
}

func cmdSize(trashMgr *trash.Manager, userUI ui.UI) {
	size, err := trashMgr.Size()
	if err != nil {
		userUI.Error(fmt.Sprintf("Failed to calculate size: %v", err))
		os.Exit(1)
	}

	userUI.Info(fmt.Sprintf("Trash size: %s", formatSize(size)))
}

func cmdConfig(cfg *config.Config, userUI ui.UI, args []string) {
	if len(args) == 0 {
		// Show all config
		fmt.Println("Current configuration:")
		fmt.Printf("  trash_dir: %s\n", cfg.TrashDir)
		fmt.Printf("  confirm_delete: %v\n", cfg.ConfirmDelete)
		fmt.Printf("  use_fancy_ui: %v\n", cfg.UseFancyUI)
		fmt.Printf("  auto_empty_days: %d\n", cfg.AutoEmptyDays)
		fmt.Printf("  max_trash_size_mb: %d\n", cfg.MaxTrashSizeMB)
		fmt.Printf("\nConfig file: %s\n", config.ConfigPath())
		return
	}

	subcommand := args[0]

	switch subcommand {
	case "get":
		if len(args) < 2 {
			userUI.Error("Usage: rc config get <key>")
			os.Exit(1)
		}
		key := args[1]
		value := cfg.Get(key)
		if value == nil {
			userUI.Error(fmt.Sprintf("Unknown config key: %s", key))
			os.Exit(1)
		}
		fmt.Printf("%s: %v\n", key, value)

	case "set":
		if len(args) < 3 {
			userUI.Error("Usage: rc config set <key> <value>")
			os.Exit(1)
		}
		key := args[1]
		value := args[2]

		// Parse value based on key
		var parsed interface{}
		switch key {
		case "trash_dir":
			parsed = value
		case "confirm_delete", "use_fancy_ui":
			parsed = value == "true" || value == "yes" || value == "1"
		case "auto_empty_days", "max_trash_size_mb":
			var intVal int
			if _, err := fmt.Sscanf(value, "%d", &intVal); err != nil {
				userUI.Error(fmt.Sprintf("Invalid integer value: %s", value))
				os.Exit(1)
			}
			parsed = intVal
		default:
			userUI.Error(fmt.Sprintf("Unknown config key: %s", key))
			os.Exit(1)
		}

		if !cfg.Set(key, parsed) {
			userUI.Error("Failed to set config value")
			os.Exit(1)
		}

		if err := cfg.Save(); err != nil {
			userUI.Error(fmt.Sprintf("Failed to save config: %v", err))
			os.Exit(1)
		}

		userUI.Success(fmt.Sprintf("Set %s = %v", key, parsed))

	case "reset":
		defaultCfg := config.DefaultConfig()
		if err := defaultCfg.Save(); err != nil {
			userUI.Error(fmt.Sprintf("Failed to reset config: %v", err))
			os.Exit(1)
		}
		userUI.Success("Config reset to defaults")

	default:
		userUI.Error(fmt.Sprintf("Unknown config subcommand: %s", subcommand))
		os.Exit(1)
	}
}

func printUsage() {
	usage := `rc - Recycle Bin Utility (GoCycled)

Usage:
  rc <command> [arguments]

Commands:
  put, trash, rm <file>...  Move files to trash
  list, ls                   List items in trash
  restore [path]             Restore item from trash (interactive if no path)
  empty                      Empty trash (permanently delete all items)
  remove <path>              Permanently delete specific item from trash
  size                       Show trash size
  config [get|set|reset]     Manage configuration
  version                    Show version
  help                       Show this help

Config Commands:
  rc config                  Show all config values
  rc config get <key>        Get a config value
  rc config set <key> <val>  Set a config value
  rc config reset            Reset config to defaults

Configuration Keys:
  trash_dir          Trash directory location
  confirm_delete     Confirm before permanent deletion (true/false)
  use_fancy_ui       Use fancy UI with Gum (true/false)
  auto_empty_days    Auto-empty trash after N days
  max_trash_size_mb  Maximum trash size in MB

Examples:
  rc put file.txt              Move file.txt to trash
  rc list                      List all trashed items
  rc restore                   Interactively restore an item
  rc restore file.txt          Restore specific file
  rc empty                     Empty trash
  rc config set use_fancy_ui true
  rc config get trash_dir

Config file: ~/.trashrc
`
	fmt.Print(usage)
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
