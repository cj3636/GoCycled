package trash

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Item represents a trashed item
type Item struct {
	OriginalPath string    `json:"original_path"`
	TrashPath    string    `json:"trash_path"`
	DeletedAt    time.Time `json:"deleted_at"`
	Size         int64     `json:"size"`
}

// Manager handles trash operations
type Manager struct {
	trashDir string
	filesDir string
	infoDir  string
}

// NewManager creates a new trash manager
func NewManager(trashDir string) (*Manager, error) {
	filesDir := filepath.Join(trashDir, "files")
	infoDir := filepath.Join(trashDir, "info")
	
	// Create trash directories if they don't exist
	if err := os.MkdirAll(filesDir, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(infoDir, 0755); err != nil {
		return nil, err
	}
	
	return &Manager{
		trashDir: trashDir,
		filesDir: filesDir,
		infoDir:  infoDir,
	}, nil
}

// Put moves a file or directory to trash
func (m *Manager) Put(path string) error {
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	
	// Check if file exists
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return err
	}
	
	// Generate unique trash filename
	baseName := filepath.Base(absPath)
	timestamp := time.Now().Format("20060102_150405")
	trashName := fmt.Sprintf("%s_%s", timestamp, baseName)
	
	trashPath := filepath.Join(m.filesDir, trashName)
	
	// Ensure unique name
	counter := 1
	for {
		if _, err := os.Stat(trashPath); os.IsNotExist(err) {
			break
		}
		trashName = fmt.Sprintf("%s_%s_%d", timestamp, baseName, counter)
		trashPath = filepath.Join(m.filesDir, trashName)
		counter++
	}
	
	// Move file to trash
	if err := os.Rename(absPath, trashPath); err != nil {
		return err
	}
	
	// Create info file
	item := Item{
		OriginalPath: absPath,
		TrashPath:    trashPath,
		DeletedAt:    time.Now(),
		Size:         getSize(fileInfo),
	}
	
	infoPath := filepath.Join(m.infoDir, trashName+".json")
	return m.saveItemInfo(infoPath, item)
}

// List returns all items in trash
func (m *Manager) List() ([]Item, error) {
	entries, err := os.ReadDir(m.infoDir)
	if err != nil {
		return nil, err
	}
	
	items := []Item{}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		
		infoPath := filepath.Join(m.infoDir, entry.Name())
		item, err := m.loadItemInfo(infoPath)
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	
	return items, nil
}

// Restore restores a file from trash to its original location
func (m *Manager) Restore(trashName string) error {
	trashPath := filepath.Join(m.filesDir, trashName)
	infoPath := filepath.Join(m.infoDir, trashName+".json")
	
	// Load item info
	item, err := m.loadItemInfo(infoPath)
	if err != nil {
		return err
	}
	
	// Check if original directory exists
	originalDir := filepath.Dir(item.OriginalPath)
	if _, err := os.Stat(originalDir); os.IsNotExist(err) {
		if err := os.MkdirAll(originalDir, 0755); err != nil {
			return err
		}
	}
	
	// Check if original path exists
	if _, err := os.Stat(item.OriginalPath); err == nil {
		return fmt.Errorf("file already exists at original location: %s", item.OriginalPath)
	}
	
	// Move file back
	if err := os.Rename(trashPath, item.OriginalPath); err != nil {
		return err
	}
	
	// Remove info file
	return os.Remove(infoPath)
}

// Remove permanently deletes an item from trash
func (m *Manager) Remove(trashName string) error {
	trashPath := filepath.Join(m.filesDir, trashName)
	infoPath := filepath.Join(m.infoDir, trashName+".json")
	
	// Remove file/directory
	if err := os.RemoveAll(trashPath); err != nil {
		return err
	}
	
	// Remove info file
	return os.Remove(infoPath)
}

// Empty removes all items from trash
func (m *Manager) Empty() error {
	// Remove all files
	if err := os.RemoveAll(m.filesDir); err != nil {
		return err
	}
	if err := os.RemoveAll(m.infoDir); err != nil {
		return err
	}
	
	// Recreate directories
	if err := os.MkdirAll(m.filesDir, 0755); err != nil {
		return err
	}
	return os.MkdirAll(m.infoDir, 0755)
}

// Size returns the total size of trash in bytes
func (m *Manager) Size() (int64, error) {
	var totalSize int64
	
	err := filepath.Walk(m.filesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	
	return totalSize, err
}

// saveItemInfo saves item info to a JSON file
func (m *Manager) saveItemInfo(path string, item Item) error {
	data, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// loadItemInfo loads item info from a JSON file
func (m *Manager) loadItemInfo(path string) (Item, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Item{}, err
	}
	
	var item Item
	if err := json.Unmarshal(data, &item); err != nil {
		return Item{}, err
	}
	
	return item, nil
}

// getSize returns the size of a file or directory
func getSize(info os.FileInfo) int64 {
	return info.Size()
}
