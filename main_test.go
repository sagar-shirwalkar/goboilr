package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// Run from root: go test ./... -v -cover
// -----------------------------------------------------------------------------

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

	// 3. Run Generator (Accessors)
	if err := GenerateFile(data, outputAccPath, false); err != nil {
		t.Fatalf("GenerateFile (Accessors) failed: %v", err)
	}
	compareFilesIgnoringWhitespace(t, outputAccPath, expectedAccPath)

	// 4. Run Generator (Constructors)
	if err := GenerateFile(data, outputConsPath, true); err != nil {
		t.Fatalf("GenerateFile (Constructors) failed: %v", err)
	}
	compareFilesIgnoringWhitespace(t, outputConsPath, expectedConsPath)
}

// compareFilesIgnoringWhitespace is a robust comparator that ignores
// indentation, newlines, and spacing differences.
func compareFilesIgnoringWhitespace(t *testing.T, generatedPath, expectedPath string) {
	genContent := normalizeTokens(t, generatedPath)
	expContent := normalizeTokens(t, expectedPath)

	if genContent != expContent {
		// If they don't match, we still print the raw files for debugging
		// but the comparison logic itself was loose.
		t.Errorf("File content mismatch (normalized).\nFile: %s\n", filepath.Base(expectedPath))
		t.Errorf("EXPECTED (Normalized Snippet): %s...", truncate(expContent, 50))
		t.Errorf("ACTUAL   (Normalized Snippet): %s...", truncate(genContent, 50))

		// Optional: Print full raw content to debug hard errors
		// t.Logf("Raw Expected:\n%s", readFile(t, expectedPath))
		// t.Logf("Raw Actual:\n%s", readFile(t, generatedPath))
	}
}

// normalizeTokens breaks the file into words/tokens and joins them with a single space.
// "func  Foo() \n {" becomes "func Foo() {"
func normalizeTokens(t *testing.T, path string) string {
	b := readFile(t, path)
	// strings.Fields splits by any whitespace (newline, tab, space)
	fields := strings.Fields(string(b))
	return strings.Join(fields, " ")
}

func readFile(t *testing.T, path string) []byte {
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", path, err)
	}
	return b
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}
