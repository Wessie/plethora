package config

import "testing"

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
