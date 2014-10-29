package config

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/Wessie/appdirs"
)

// this file handles locating the database, the database is a directory with
// multiple database files in it.

// maxLocLength is the maximum length of the database path
const maxLocLength = 4096

// dbloc keeps track of the database location and keeps an open file handle
// to lock against.
var dbloc struct {
	// loc is the current location of the database
	loc string
	// open handle to the dbloc file
	*os.File
	// mutex protects the other fields
	sync.Mutex
}

var (
	dir = appdirs.UserConfigDir(Name, "", "", false)
	// dblocFile is the location of the dbloc file, the dbloc
	// file contains a single path. This path is the location
	// of the database and should be a directory
	dblocFile = func() string {
		return filepath.Join(dir, "dbloc")
	}
	// defaultLoc is the default database location
	defaultLoc = func() string {
		return filepath.Join(dir, "db")
	}
)

// Init initializes the config package, this should be called before any other
// functions are used in this package.
func Init() error {
	// TODO: lock file for concurrent access
	f, err := os.OpenFile(dblocFile(), os.O_CREATE|os.O_RDWR, 0770)
	if err != nil {
		return err
	}
	dbloc.File = f

	b, err := ioutil.ReadAll(&io.LimitedReader{R: f, N: maxLocLength})
	if err != nil {
		return err
	}

	if len(b) > 0 {
		dbloc.loc = string(b)
		return nil
	}

	return UpdateLocation(defaultLoc())
}

// Location returns the current location of the database
func Location() string {
	dbloc.Lock()
	defer dbloc.Unlock()
	return dbloc.loc
}

// UpdateLocation updates the database location in the dbloc file, this
// does not move any existing databases around.
func UpdateLocation(loc string) error {
	if len(loc) > maxLocLength {
		return errors.New("new database location is too long")
	}

	dbloc.Lock()
	defer dbloc.Unlock()

	_, err := dbloc.WriteAt([]byte(loc), 0)
	if err != nil {
		return err
	}

	dbloc.loc = loc
	return nil
}
