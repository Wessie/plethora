package config

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func createTempDir() (string, error) {
	hex := fmt.Sprintf("%x", rand.Int())

	d := filepath.Join(os.TempDir(), hex)
	if err := os.MkdirAll(d, 0770); err != nil {
		return "", err
	}

	return d, nil
}

func removeTempDir(dir string) {
	os.RemoveAll(dir)
}

func TestInit(t *testing.T) {
	var err error
	// change from the user-looking directory to a temp one
	dir, err = createTempDir()
	if err != nil {
		t.Fatal("could not create temporary directory")
	}
	defer removeTempDir(dir)

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
