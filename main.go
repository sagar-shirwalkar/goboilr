package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// We pass real OS arguments and the real os.Getenv function
	// os.Args[1:] skips the program name itself
	if err := run(os.Args[1:], os.Getenv); err != nil {
		log.Fatal(err)
	}
}

// run contains the main logic, decoupled from os.Exit for testing
func run(args []string, getEnv func(string) string) error {
	// 1. Determine Input File
	// Use a local FlagSet to avoid polluting global state during tests
	fs := flag.NewFlagSet("goboilr", flag.ContinueOnError)
	fileFlag := fs.String("file", "", "The Go source file to process")

	// Parse the provided arguments
	if err := fs.Parse(args); err != nil {
		return err
	}

	targetFile := *fileFlag
	if targetFile == "" {
		targetFile = getEnv("GOFILE")
	}

	if targetFile == "" {
		return fmt.Errorf("error: no file specified. Use -file flag or run via 'go generate'")
	}

	// 2. Resolve Absolute Paths
	absPath, err := filepath.Abs(targetFile)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	dir := filepath.Dir(absPath)
	filename := filepath.Base(absPath)
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)

	// 3. Define Output Names
	accessorsFile := filepath.Join(dir, fmt.Sprintf("%s_accessors.go", baseName))
	constructorsFile := filepath.Join(dir, fmt.Sprintf("%s_constructors.go", baseName))

	fmt.Printf("GoBoilr: Processing %s\n", filename)

	// 4. Parse & Generate
	data, err := ParseFile(absPath)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	if err := GenerateFile(data, accessorsFile, false); err != nil {
		return fmt.Errorf("failed to generate accessors: %w", err)
	}
	fmt.Printf("  -> Created %s\n", filepath.Base(accessorsFile))

	if err := GenerateFile(data, constructorsFile, true); err != nil {
		return fmt.Errorf("failed to generate constructors: %w", err)
	}
	fmt.Printf("  -> Created %s\n", filepath.Base(constructorsFile))

	return nil
}
