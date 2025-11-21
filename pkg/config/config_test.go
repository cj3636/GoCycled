package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.TrashDir == "" {
		t.Error("TrashDir should not be empty")
	}

	if !cfg.ConfirmDelete {
		t.Error("ConfirmDelete should be true by default")
	}

	if cfg.UseFancyUI {
		t.Error("UseFancyUI should be false by default")
	}

	if cfg.AutoEmptyDays != 30 {
		t.Errorf("AutoEmptyDays should be 30, got %d", cfg.AutoEmptyDays)
	}

	if cfg.MaxTrashSizeMB != 1024 {
		t.Errorf("MaxTrashSizeMB should be 1024, got %d", cfg.MaxTrashSizeMB)
	}
}

func TestConfigSaveLoad(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()

	// Save original ConfigPath function behavior
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create and save config
	cfg := DefaultConfig()
	cfg.UseFancyUI = true
	cfg.AutoEmptyDays = 45

	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify loaded values
	if !loadedCfg.UseFancyUI {
		t.Error("UseFancyUI should be true")
	}

	if loadedCfg.AutoEmptyDays != 45 {
		t.Errorf("AutoEmptyDays should be 45, got %d", loadedCfg.AutoEmptyDays)
	}
}

func TestConfigGetSet(t *testing.T) {
	cfg := DefaultConfig()

	// Test Get
	if val := cfg.Get("trash_dir"); val == nil {
		t.Error("Get should return trash_dir")
	}

	if val := cfg.Get("confirm_delete"); val != true {
		t.Error("Get should return true for confirm_delete")
	}

	if val := cfg.Get("invalid_key"); val != nil {
		t.Error("Get should return nil for invalid key")
	}

	// Test Set
	if !cfg.Set("confirm_delete", false) {
		t.Error("Set should succeed for confirm_delete")
	}

	if cfg.ConfirmDelete {
		t.Error("ConfirmDelete should be false after Set")
	}

	if !cfg.Set("auto_empty_days", 60) {
		t.Error("Set should succeed for auto_empty_days")
	}

	if cfg.AutoEmptyDays != 60 {
		t.Error("AutoEmptyDays should be 60 after Set")
	}

	if cfg.Set("invalid_key", "value") {
		t.Error("Set should fail for invalid key")
	}
}
