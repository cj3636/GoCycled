package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	TrashDir        string `json:"trash_dir"`
	ConfirmDelete   bool   `json:"confirm_delete"`
	UseFancyUI      bool   `json:"use_fancy_ui"`
	AutoEmptyDays   int    `json:"auto_empty_days"`
	MaxTrashSizeMB  int    `json:"max_trash_size_mb"`
}

// DefaultConfig returns a new Config with default values
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		TrashDir:        filepath.Join(homeDir, ".local", "share", "Trash"),
		ConfirmDelete:   true,
		UseFancyUI:      false,
		AutoEmptyDays:   30,
		MaxTrashSizeMB:  1024,
	}
}

// ConfigPath returns the path to the config file
func ConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".trashrc")
}

// Load loads the config from ~/.trashrc or creates a default one
func Load() (*Config, error) {
	configPath := ConfigPath()
	
	// If config doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := DefaultConfig()
		if err := cfg.Save(); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	
	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	
	return cfg, nil
}

// Save saves the config to ~/.trashrc
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	
	configPath := ConfigPath()
	return os.WriteFile(configPath, data, 0644)
}

// Get retrieves a config value by key
func (c *Config) Get(key string) interface{} {
	switch key {
	case "trash_dir":
		return c.TrashDir
	case "confirm_delete":
		return c.ConfirmDelete
	case "use_fancy_ui":
		return c.UseFancyUI
	case "auto_empty_days":
		return c.AutoEmptyDays
	case "max_trash_size_mb":
		return c.MaxTrashSizeMB
	default:
		return nil
	}
}

// Set sets a config value by key
func (c *Config) Set(key string, value interface{}) bool {
	switch key {
	case "trash_dir":
		if v, ok := value.(string); ok {
			c.TrashDir = v
			return true
		}
	case "confirm_delete":
		if v, ok := value.(bool); ok {
			c.ConfirmDelete = v
			return true
		}
	case "use_fancy_ui":
		if v, ok := value.(bool); ok {
			c.UseFancyUI = v
			return true
		}
	case "auto_empty_days":
		if v, ok := value.(int); ok {
			c.AutoEmptyDays = v
			return true
		}
	case "max_trash_size_mb":
		if v, ok := value.(int); ok {
			c.MaxTrashSizeMB = v
			return true
		}
	}
	return false
}
