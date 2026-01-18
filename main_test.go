package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// Unit Tests (Helper Functions)
// -----------------------------------------------------------------------------

func TestHelpers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		function func(string) string
	}{
		// Capitalize
		{"Cap Normal", "hello", "Hello", capitalize},
		{"Cap Empty", "", "", capitalize},
		{"Cap 1 char", "h", "H", capitalize},

		// LowerFirst
		{"Lower Normal", "Hello", "hello", lowerFirst},
		{"Lower Empty", "", "", lowerFirst},
		{"Lower 1 char", "H", "h", lowerFirst},

		// ExtractEmbeddedName
		{"Embed Simple", "Person", "Person", extractEmbeddedName},
		{"Embed Pointer", "*Person", "Person", extractEmbeddedName},
		{"Embed Package", "models.Person", "Person", extractEmbeddedName},
		{"Embed Ptr Pkg", "*models.Person", "Person", extractEmbeddedName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.function(tt.input)
			if got != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, got)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// Integration Tests (Parser & Generator)
// -----------------------------------------------------------------------------

func TestParserErrors(t *testing.T) {
	// Test non-existent file
	_, err := ParseFile("non_existent_file.go")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestGeneratorErrors(t *testing.T) {
	// Test writing to an invalid path
	data := &FileData{PackageName: "test"}
	err := GenerateFile(data, "/invalid/path/output.go", false)
	if err == nil {
		t.Error("Expected error writing to invalid path, got nil")
	}
}

func TestEndToEnd(t *testing.T) {
	// 1. Setup paths
	wd, _ := os.Getwd()
	sourcePath := filepath.Join(wd, "testdata", "source.go")
	expectedAccPath := filepath.Join(wd, "testdata", "expected_acc.go")
	expectedConsPath := filepath.Join(wd, "testdata", "expected_cons.go")

	tmpDir := t.TempDir()
	outputAccPath := filepath.Join(tmpDir, "output_accessors.go")
	outputConsPath := filepath.Join(tmpDir, "output_constructors.go")

	// 2. Run Parser
	data, err := ParseFile(sourcePath)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Basic Validation
	if data.PackageName != "testdata" {
		t.Errorf("Expected package 'testdata', got '%s'", data.PackageName)
	}
	// We expect 2 structs: ComplexStruct and Base
	if len(data.Structs) != 2 {
		t.Fatalf("Expected 2 structs, got %d", len(data.Structs))
	}

	// 3. Run Generator (Accessors)
	if err := GenerateFile(data, outputAccPath, false); err != nil {
		t.Fatalf("GenerateFile (Accessors) failed: %v", err)
	}
	compareFiles(t, outputAccPath, expectedAccPath)

	// 4. Run Generator (Constructors)
	if err := GenerateFile(data, outputConsPath, true); err != nil {
		t.Fatalf("GenerateFile (Constructors) failed: %v", err)
	}
	compareFiles(t, outputConsPath, expectedConsPath)
}

func compareFiles(t *testing.T, generatedPath, expectedPath string) {
	genContent := normalize(t, generatedPath)
	expContent := normalize(t, expectedPath)

	if genContent != expContent {
		t.Errorf("File content mismatch for %s.\n\nEXPECTED:\n%s\n\nACTUAL:\n%s",
			filepath.Base(expectedPath), expContent, genContent)
	}
}

func normalize(t *testing.T, path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", path, err)
	}
	s := string(b)
	// Normalize Windows/Linux line endings
	s = strings.ReplaceAll(s, "\r\n", "\n")
	// Normalize leading/trailing whitespace (fixes EOF mismatch)
	return strings.TrimSpace(s)
}
