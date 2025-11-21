package trash

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewManager(t *testing.T) {
	tempDir := t.TempDir()

	mgr, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	if mgr.trashDir != tempDir {
		t.Errorf("Expected trashDir %s, got %s", tempDir, mgr.trashDir)
	}

	// Check if directories were created
	if _, err := os.Stat(mgr.filesDir); os.IsNotExist(err) {
		t.Error("Files directory should be created")
	}

	if _, err := os.Stat(mgr.infoDir); os.IsNotExist(err) {
		t.Error("Info directory should be created")
	}
}

func TestPutAndList(t *testing.T) {
	tempDir := t.TempDir()
	mgr, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Create test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Put file in trash
	if err := mgr.Put(testFile); err != nil {
		t.Fatalf("Failed to put file in trash: %v", err)
	}

	// Verify original file is gone
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("Original file should be deleted")
	}

	// List items
	items, err := mgr.List()
	if err != nil {
		t.Fatalf("Failed to list items: %v", err)
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}

	if len(items) > 0 {
		if items[0].OriginalPath != testFile {
			t.Errorf("Expected original path %s, got %s", testFile, items[0].OriginalPath)
		}
	}
}

func TestRestore(t *testing.T) {
	tempDir := t.TempDir()
	mgr, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Create test file
	testFile := filepath.Join(tempDir, "test_restore.txt")
	testContent := []byte("restore test")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Put file in trash
	if err := mgr.Put(testFile); err != nil {
		t.Fatalf("Failed to put file in trash: %v", err)
	}

	// Get trash name
	items, _ := mgr.List()
	if len(items) == 0 {
		t.Fatal("No items in trash")
	}
	trashName := filepath.Base(items[0].TrashPath)

	// Restore file
	if err := mgr.Restore(trashName); err != nil {
		t.Fatalf("Failed to restore file: %v", err)
	}

	// Verify file is restored
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Error("Restored file should exist")
	}

	if string(content) != string(testContent) {
		t.Error("Restored file content should match original")
	}

	// Verify trash is empty
	items, _ = mgr.List()
	if len(items) != 0 {
		t.Error("Trash should be empty after restore")
	}
}

func TestRemove(t *testing.T) {
	tempDir := t.TempDir()
	mgr, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Create test file
	testFile := filepath.Join(tempDir, "test_remove.txt")
	if err := os.WriteFile(testFile, []byte("remove test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Put file in trash
	if err := mgr.Put(testFile); err != nil {
		t.Fatalf("Failed to put file in trash: %v", err)
	}

	// Get trash name
	items, _ := mgr.List()
	if len(items) == 0 {
		t.Fatal("No items in trash")
	}
	trashName := filepath.Base(items[0].TrashPath)

	// Remove file
	if err := mgr.Remove(trashName); err != nil {
		t.Fatalf("Failed to remove file: %v", err)
	}

	// Verify trash is empty
	items, _ = mgr.List()
	if len(items) != 0 {
		t.Error("Trash should be empty after remove")
	}

	// Verify original file is still gone
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("Original file should not be restored")
	}
}

func TestEmpty(t *testing.T) {
	tempDir := t.TempDir()
	mgr, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Create and trash multiple files
	for i := 0; i < 3; i++ {
		testFile := filepath.Join(tempDir, filepath.Base(tempDir)+".txt")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		mgr.Put(testFile)
	}

	// Empty trash
	if err := mgr.Empty(); err != nil {
		t.Fatalf("Failed to empty trash: %v", err)
	}

	// Verify trash is empty
	items, _ := mgr.List()
	if len(items) != 0 {
		t.Errorf("Trash should be empty, got %d items", len(items))
	}
}

func TestSize(t *testing.T) {
	tempDir := t.TempDir()
	mgr, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Create test file
	testContent := []byte("test content for size")
	testFile := filepath.Join(tempDir, "test_size.txt")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Put file in trash
	if err := mgr.Put(testFile); err != nil {
		t.Fatalf("Failed to put file in trash: %v", err)
	}

	// Get size
	size, err := mgr.Size()
	if err != nil {
		t.Fatalf("Failed to get size: %v", err)
	}

	if size == 0 {
		t.Error("Size should not be 0")
	}
}
