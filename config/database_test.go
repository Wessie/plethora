package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestDatabaseClosing(t *testing.T) {
	defer TestConfiguration()()

	d, err := Database("closing")
	if err != nil {
		t.Fatal("unable to open database:", err)
	}

	d.Close()

	d, err = Database("closing")
	if err != nil {
		t.Fatal("unable to reopen database:", err)
	}

	d.Close()
}

func TestDatabaseCache(t *testing.T) {
	defer TestConfiguration()()

	d1, err := Database("cache")
	if err != nil {
		t.Fatal("unable to open first database:", err)
	}

	d2, err := Database("cache")
	if err != nil {
		t.Fatal("unable to open second database:", err)
	}

	if d1 != d2 {
		t.Fatal("received two different databases with same name")
	}
}

func TestDatabaseName(t *testing.T) {
	defer TestConfiguration()()

	a, b := DatabasePath("dup"), DatabasePath("dup")
	if a != b {
		t.Errorf("database filepath is not deterministic: %s != %s", a, b)
	}

	path := DatabasePath("nametest")
	if !filepath.IsAbs(path) {
		t.Error("database filepath is not absolute:", path)
	}

	if strings.Count(filepath.ToSlash(path), "/") < 2 {
		// this is possible to hit legimately if XDG* environment variables
		// are set to the root of the system, but that is a rare occurence
		// and instead we assume it is wrong for that to happen.
		t.Error("database filepath is too close to the root:", path)
	}
}
