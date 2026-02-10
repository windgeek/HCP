package manifest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"github.com/windgeek/HCP/pkg/identity"
)

func TestManifestSigning(t *testing.T) {
	// 1. Setup Keys
	key, err := identity.GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	// 2. Setup Test Content
	tmpDir := t.TempDir()
	contentPath := filepath.Join(tmpDir, "test_content.txt")
	err = os.WriteFile(contentPath, []byte("Hello HCP"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test content: %v", err)
	}

	// 3. Create Manifest
	authorAddr := "bc1qtest..." // Mock address
	m, err := NewManifest(contentPath, authorAddr)
	if err != nil {
		t.Fatalf("NewManifest failed: %v", err)
	}

	if m.ContentHash == "" {
		t.Fatal("ContentHash is empty")
	}

	// 4. Sign
	err = m.Sign(key)
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	if m.Signature == "" {
		t.Fatal("Signature is empty")
	}

	// 5. Save and Reload
	manifestPath := contentPath + ".hcp"
	err = m.Save(manifestPath)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("Failed to read manifest file: %v", err)
	}

	var loadedM Manifest
	err = json.Unmarshal(data, &loadedM)
	if err != nil {
		t.Fatalf("Failed to unmarshal manifest: %v", err)
	}

	if loadedM.Signature != m.Signature {
		t.Fatal("Loaded signature does not match")
	}
}
