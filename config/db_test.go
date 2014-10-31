package config

import "testing"

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

func TestDatabaseName(t *testing.T) {}
