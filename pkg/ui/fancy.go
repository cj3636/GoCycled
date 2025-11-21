package ui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cj3636/GoCycled/pkg/trash"
)

// FancyUI is a Gum-based fancy UI
type FancyUI struct {
	basic *BasicUI
}

// NewFancyUI creates a new fancy UI
func NewFancyUI() *FancyUI {
	return &FancyUI{
		basic: NewBasicUI(),
	}
}

// Confirm asks for user confirmation using Gum
func (u *FancyUI) Confirm(message string) bool {
	// Check if gum is available
	if !isGumAvailable() {
		return u.basic.Confirm(message)
	}
	
	cmd := exec.Command("gum", "confirm", message)
	err := cmd.Run()
	return err == nil
}

// DisplayItems displays a list of trash items
func (u *FancyUI) DisplayItems(items []trash.Item) {
	if len(items) == 0 {
		u.Info("Trash is empty")
		return
	}
	
	// Use basic display for now
	u.basic.DisplayItems(items)
}

// SelectItem prompts user to select an item using Gum
func (u *FancyUI) SelectItem(items []trash.Item) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to select")
	}
	
	// Check if gum is available
	if !isGumAvailable() {
		return u.basic.SelectItem(items)
	}
	
	// Build options
	options := make([]string, len(items))
	for i, item := range items {
		options[i] = item.OriginalPath
	}
	
	cmd := exec.Command("gum", "choose", "--header=Select item to restore:")
	cmd.Stdin = strings.NewReader(strings.Join(options, "\n"))
	
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("selection cancelled")
	}
	
	selected := strings.TrimSpace(string(output))
	
	// Find matching item
	for _, item := range items {
		if item.OriginalPath == selected {
			return item.TrashPath, nil
		}
	}
	
	return "", fmt.Errorf("item not found")
}

// Success displays a success message
func (u *FancyUI) Success(message string) {
	if isGumAvailable() {
		cmd := exec.Command("gum", "style", "--foreground=2", "✓ "+message)
		cmd.Run()
	} else {
		u.basic.Success(message)
	}
}

// Error displays an error message
func (u *FancyUI) Error(message string) {
	if isGumAvailable() {
		cmd := exec.Command("gum", "style", "--foreground=1", "✗ Error: "+message)
		cmd.Run()
	} else {
		u.basic.Error(message)
	}
}

// Info displays an info message
func (u *FancyUI) Info(message string) {
	if isGumAvailable() {
		cmd := exec.Command("gum", "style", "--foreground=4", "ℹ "+message)
		cmd.Run()
	} else {
		u.basic.Info(message)
	}
}

// isGumAvailable checks if gum command is available
func isGumAvailable() bool {
	_, err := exec.LookPath("gum")
	return err == nil
}
