package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestReadFiles(t *testing.T) {
	// Prepare test files
	file1 := "testfile1.txt"
	file2 := "testfile2.txt"
	file1Content := "Test file 1 content."
	file2Content := "Test file 2 content."

	err := os.WriteFile(file1, []byte(file1Content), 0644)
	if err != nil {
		log.Fatalf("Error creating test file %s: %s", file1, err)
	}
	defer os.Remove(file1)

	err = os.WriteFile(file2, []byte(file2Content), 0644)
	if err != nil {
		log.Fatalf("Error creating test file %s: %s", file2, err)
	}
	defer os.Remove(file2)

	// Test readFiles function
	expected := strings.Builder{}
	expected.WriteString(fmt.Sprintf("PRINT %s\n", file1))
	expected.WriteString("FILE_START\n")
	expected.WriteString(file1Content)
	expected.WriteString("\nFILE_END\n\n")
	expected.WriteString(fmt.Sprintf("PRINT %s\n", file2))
	expected.WriteString("FILE_START\n")
	expected.WriteString(file2Content)
	expected.WriteString("\nFILE_END\n\n")

	result := readFiles([]string{file1, file2})

	if result != expected.String() {
		t.Errorf("readFiles failed, expected:\n%s, but got:\n%s", expected.String(), result)
	}
}
