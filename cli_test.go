package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	// Setup a temporary directory for file operations
	tmpDir := t.TempDir()

	// Create a dummy Go file to process
	dummySource := filepath.Join(tmpDir, "dummy.go")
	dummyContent := `package dummy
	// gen:new
	type User struct { 
		Name string ` + "`gen:\"get\"`" + `
	}`
	if err := os.WriteFile(dummySource, []byte(dummyContent), 0644); err != nil {
		t.Fatalf("Failed to create dummy file: %v", err)
	}

	tests := []struct {
		name        string
		args        []string
		mockEnv     map[string]string
		expectError bool
	}{
		{
			name:        "Error: No arguments",
			args:        []string{},
			mockEnv:     nil,
			expectError: true,
		},
		{
			name:        "Success: Via Flag",
			args:        []string{"-file", dummySource},
			mockEnv:     nil,
			expectError: false,
		},
		{
			name:        "Success: Via Env Var",
			args:        []string{},
			mockEnv:     map[string]string{"GOFILE": dummySource},
			expectError: false,
		},
		{
			name:        "Error: Invalid File Path",
			args:        []string{"-file", filepath.Join(tmpDir, "non_existent.go")},
			mockEnv:     nil,
			expectError: true,
		},
		{
			name:        "Error: Invalid Flags",
			args:        []string{"-unknownFlag"},
			mockEnv:     nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock Env function
			mockGetEnv := func(key string) string {
				if tt.mockEnv != nil {
					return tt.mockEnv[key]
				}
				return ""
			}

			err := run(tt.args, mockGetEnv)

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// If success, verify files were created
			if !tt.expectError {
				accFile := filepath.Join(tmpDir, "dummy_accessors.go")
				if _, err := os.Stat(accFile); os.IsNotExist(err) {
					t.Errorf("Expected output file %s was not created", accFile)
				}
			}
		})
	}
}
