package config

import (
	"path/filepath"
	"testing"

	"github.com/Wessie/plethora/config/testutil"
)

func TestConfigInit(t *testing.T) {
	// change from the user-looking directory to a temp one
	dir, err := testutil.TempDir()
	if err != nil {
		t.Fatal("could not create temporary directory")
	}
	defer testutil.RemoveTempDir(dir)

	dbLocFile := filepath.Join(dir, "dbloc")

	dbloc := dbLoc{
		filename: dbLocFile,
	}

	// test if initializing works
	if err := dbloc.init(); err != nil {
		t.Fatal("failed initializing:", err)
	}

	if dbloc.path == "" {
		t.Error("failed: set default location")
	}

	newLoc := "testing"
	if err := dbloc.updateFile(newLoc); err != nil {
		t.Error("failed: update location", err)
	}

	if dbloc.path != newLoc {
		t.Error("failed: location is not updated")
	}
}
