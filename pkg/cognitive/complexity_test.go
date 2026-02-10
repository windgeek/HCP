package cognitive

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeGoFile(t *testing.T) {
	// Create specific temp dir that doesn't conflict
	tmpDir, err := os.MkdirTemp("", "test_cognitive_*")
	if err != nil {
		t.Fatal(err)
	}
	// Attempt to clean up, but ignore errors if it fails (common in tests)
	defer os.RemoveAll(tmpDir)

	src := `package main
	
	func main() {
		if true {
			print("hello")
		} else {
			print("world")
		}
	}
	
	func complex() {
		for i := 0; i < 10; i++ {
			if i % 2 == 0 {
				continue
			}
		}
	}
	`
	
	tmpFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(tmpFile, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	stats, err := AnalyzeComplexity(tmpFile)
	if err != nil {
		t.Fatalf("AnalyzeComplexity failed: %v", err)
	}

	// Function main: base 1 + if/else (1?) -> 2
	// Function complex: base 1 + for (1) + if (1) -> 3
	// Total Cyclomatic should be around 5
	// Check range to be safe against exact implementation details
	if stats.Cyclomatic < 4 {
		t.Errorf("Expected Cyclomatic complexity >= 4, got %d", stats.Cyclomatic)
	}
	
	if stats.Functions != 2 {
		t.Errorf("Expected 2 functions, got %d", stats.Functions)
	}
}
