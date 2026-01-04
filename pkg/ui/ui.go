package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cj3636/GoCycled/pkg/trash"
)

// UI interface for different UI implementations
type UI interface {
	Confirm(message string) bool
	DisplayItems(items []trash.Item)
	SelectItem(items []trash.Item) (string, error)
	Success(message string)
	Error(message string)
	Info(message string)
}

// BasicUI is a simple text-based UI
type BasicUI struct {
	reader *bufio.Reader
}

// NewBasicUI creates a new basic UI
func NewBasicUI() *BasicUI {
	return &BasicUI{
		reader: bufio.NewReader(os.Stdin),
	}
}

// Confirm asks for user confirmation
func (u *BasicUI) Confirm(message string) bool {
	fmt.Printf("%s (y/N): ", message)
	response, _ := u.reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// DisplayItems displays a list of trash items
func (u *BasicUI) DisplayItems(items []trash.Item) {
	if len(items) == 0 {
		fmt.Println("Trash is empty")
		return
	}

	fmt.Printf("\n%-40s %-20s %-15s\n", "Original Path", "Deleted At", "Size")
	fmt.Println(strings.Repeat("-", 80))

	for _, item := range items {
		size := formatSize(item.Size)
		deletedAt := item.DeletedAt.Format("2006-01-02 15:04:05")
		fmt.Printf("%-40s %-20s %-15s\n",
			truncate(item.OriginalPath, 40),
			deletedAt,
			size)
	}
	fmt.Println()
}

// SelectItem prompts user to select an item
func (u *BasicUI) SelectItem(items []trash.Item) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to select")
	}

	for i, item := range items {
		fmt.Printf("%d. %s\n", i+1, item.OriginalPath)
	}

	fmt.Print("\nSelect item number: ")
	input, err := u.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %v", err)
	}

	var selection int
	input = strings.TrimSpace(input)
	if _, err := fmt.Sscanf(input, "%d", &selection); err != nil || selection < 1 || selection > len(items) {
		return "", fmt.Errorf("invalid selection")
	}

	return items[selection-1].TrashPath, nil
}

// Success displays a success message
func (u *BasicUI) Success(message string) {
	fmt.Printf("✓ %s\n", message)
}

// Error displays an error message
func (u *BasicUI) Error(message string) {
	fmt.Fprintf(os.Stderr, "✗ Error: %s\n", message)
}

// Info displays an info message
func (u *BasicUI) Info(message string) {
	fmt.Printf("ℹ %s\n", message)
}

// Helper functions

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

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
