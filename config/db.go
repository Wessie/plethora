package config

import (
	"path/filepath"
	"sync"

	"github.com/boltdb/bolt"
)

var dbMutex sync.Mutex
var dbMap map[string]*bolt.DB

// Database returns the bolt database associated with the name
// given. If Database is called multiple times with the same name
// it will return the same *bolt.DB unless it was previously closed.
func Database(name string) (*bolt.DB, error) {
	path := DatabasePath(name)

	dbMutex.Lock()
	defer dbMutex.Unlock()

	if db, ok := dbMap[path]; ok {
		tx, err := db.Begin(false)
		// if everything is fine, return what we have
		if err == nil {
			tx.Rollback()
			return db, nil
		}

		// if we otherwise got some unexpected error, ignore
		// it and have the caller deal with it when they open
		// a transaction.
		if err != bolt.ErrDatabaseNotOpen {
			return db, nil
		}
	}

	// if we either have no existing db, or have one that is closed
	// open a new one instead
	db, err := bolt.Open(path, 0660, nil)
	if err != nil {
		return nil, err
	}

	dbMap[path] = db
	return db, nil
}

// DatabasePath returns the filepath to the database associated with
// the name given. This function can be used to use an alternate file-based
// database in code.
func DatabasePath(name string) string {
	return filepath.Join(Location(), name)
}
