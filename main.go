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
	// 1. Determine Input File
	// Check flags first, then environment variable (set by go generate)
	fileFlag := flag.String("file", "", "The Go source file to process")
	flag.Parse()

	targetFile := *fileFlag
	if targetFile == "" {
		targetFile = os.Getenv("GOFILE")
	}

	if targetFile == "" {
		log.Fatal("Error: No file specified. Use -file flag or run via 'go generate'")
	}

	// 2. Resolve Absolute Paths
	// This ensures logic works whether we are in root or subfolder
	absPath, err := filepath.Abs(targetFile)
	if err != nil {
		log.Fatal(err)
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
		log.Fatalf("Parsing failed: %v", err)
	}

	if err := GenerateFile(data, accessorsFile, false); err != nil {
		log.Fatalf("Failed to generate accessors: %v", err)
	}
	fmt.Printf("  -> Created %s\n", filepath.Base(accessorsFile))

	if err := GenerateFile(data, constructorsFile, true); err != nil {
		log.Fatalf("Failed to generate constructors: %v", err)
	}
	fmt.Printf("  -> Created %s\n", filepath.Base(constructorsFile))
}
