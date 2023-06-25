package main

import (
	"os"
	"strings"
	"testing"
)

func TestUpdateFile(t *testing.T) {
	// create a file
	tmpfile, err := os.CreateTemp("", "test_update_file")
	if err != nil {
		t.Fatalf("unable to create tempfile: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// update the file
	contents := []string{"Line 1", "Line 2", "Line 3"}
	err = updateFile(tmpfile.Name(), contents)
	if err != nil {
		t.Fatalf("unable to update file: %v", err)
	}

	// check that the file was updated
	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("unable to read updated file: %v", err)
	}
	expected := "Line 1\nLine 2\nLine 3"
	if string(data) != expected {
		t.Errorf("expected file content to be %q but got %q", expected, string(data))
	}
}

func TestDeleteFile(t *testing.T) {
	// create a file
	tmpfile, err := os.CreateTemp("", "test_delete_file")
	if err != nil {
		t.Fatalf("unable to create tempfile: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// delete the file
	err = deleteFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("unable to delete file: %v", err)
	}

	// check that the file was deleted
	_, err = os.Stat(tmpfile.Name())
	if err == nil || !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("expected file %s to not exist, but it does", tmpfile.Name())
	}
}

func TestApplyChanges(t *testing.T) {
	// create a file and update a file
	tmpfile1, err := os.CreateTemp("", "test_create_file")
	if err != nil {
		t.Fatalf("unable to create tempfile: %v", err)
	}
	defer os.Remove(tmpfile1.Name())
	tmpfile2, err := os.CreateTemp("", "test_create_file2")
	if err != nil {
		t.Fatalf("unable to create tempfile: %v", err)
	}
	defer os.Remove(tmpfile2.Name())

	changes := "CREATE " + tmpfile1.Name() + "\nFILE_START\ncreated file contents...\nFILE_END\n" +
		"UPDATE " + tmpfile1.Name() + "\nFILE_START\nupdated file contents...\nFILE_END\n" +
		"DELETE " + tmpfile2.Name() + "\nFILE_START\n(deleted)\nFILE_END"
	t.Setenv("APPLY_CHANGES_NO_CONFIRM", "true")
	err = applyChanges(changes)
	if err != nil {
		t.Fatalf("unable to apply changes: %v", err)
	}

	// check that tmpfile1 was created and updated
	data, err := os.ReadFile(tmpfile1.Name())
	if err != nil {
		t.Fatalf("unable to read updated file: %v", err)
	}

	expected := "updated file contents..."
	if string(data) != expected {
		t.Errorf("expected file content to be %q but got %q", expected, string(data))
	}

	// check that tmpfile2 was deleted
	_, err = os.Stat(tmpfile2.Name())
	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("expected file %s to not exist, but it does", tmpfile2.Name())
	}
}
