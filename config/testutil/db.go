package testutil

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
)

func NewTestDB() (*bolt.DB, func(), error) {
	dir, err := TempDir()
	if err != nil {
		return nil, func() {}, err
	}

	// create a random name for our database
	name := fmt.Sprintf("%x", rand.Int())
	dbFile := filepath.Join(dir, name)

	db, err := bolt.Open(dbFile, 0770, nil)

	// function to cleanup all the things we've created here, this
	// needs to be called by the caller when results are no longer
	// needed.
	cleanup := func() {
		db.Close()
		RemoveTempDir(dir)
	}

	return db, cleanup, err
}

func TempDir() (string, error) {
	hex := fmt.Sprintf("%x", rand.Int())

	d := filepath.Join(os.TempDir(), hex)
	if err := os.MkdirAll(d, 0770); err != nil {
		return "", err
	}

	return d, nil
}

func RemoveTempDir(dir string) {
	os.RemoveAll(dir)
}
