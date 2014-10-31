package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigInit(t *testing.T) {
	defer TestConfiguration()()

	// test if initializing works
	if err := Init(); err != nil {
		t.Fatal("failed initializing:", err)
	}

	if Location() == "" {
		t.Error("failed: set default location")
	}

	newLoc := "testing"
	if err := UpdateLocation(newLoc); err != nil {
		t.Error("failed: update location", err)
	}

	if Location() != newLoc {
		t.Error("failed: location is not updated")
	}
}

func TestTestingConfiguration(t *testing.T) {
	loc := dbLocation

	cleanup := TestConfiguration()

	if loc.filename == dbLocation.filename {
		t.Error("test location is empty string, instead of temporary filepath")
	}

	if loc.path == dbLocation.path {
		t.Error("test database location is empty string, instead of temporary directory")
	}

	f, err := os.Create(dbLocation.filename)
	if err != nil {
		t.Fatal("failed to create file in temporary directory", err)
	}
	f.Close()

	// keep a copy around for after cleanup
	filename := dbLocation.filename
	// check if the cleanup function properly deletes everything
	cleanup()

	_, err = os.Stat(filename)
	if !os.IsNotExist(err) {
		t.Error("cleanup did not remove temporary file")
	}

	_, err = os.Stat(filepath.Dir(filename))
	if !os.IsNotExist(err) {
		t.Error("cleanup did not remove temporary directory")
	}
}
