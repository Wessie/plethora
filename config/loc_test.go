package config

import (
	"testing"

	"github.com/Wessie/plethora/config/testutil"
)

func TestConfigInit(t *testing.T) {
	var err error
	// change from the user-looking directory to a temp one
	dir, err = testutil.TempDir()
	if err != nil {
		t.Fatal("could not create temporary directory")
	}
	defer testutil.RemoveTempDir(dir)

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
