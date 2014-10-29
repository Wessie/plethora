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

var (
	// dbLocation is the global tracker of the location file
	dbLocation dbLoc

	// configDir is the directory we look in for the dbloc file
	configDir = appdirs.UserConfigDir(Name, "", "", false)

	// dbLocationFile is the filepath to the dbloc file
	dbLocationFile = filepath.Join(configDir, "dbloc")

	// dbDefaultLocation is the default path to the database files
	dbDefaultLocation = filepath.Join(configDir, "db")
)

// dbLoc is the location of the database files, stored inside the file located
// at the filepath returned by dbLocationFile
type dbLoc struct {
	// path is the current location of the database
	path string
	// mutex protects the path field
	sync.Mutex
}

// openFile opens and returns the dbLocationFile with correct flags
// and permissions set.
func (d *dbLoc) openFile() (*os.File, error) {
	// TODO: lock file for concurrent access
	return os.OpenFile(dbLocationFile, os.O_CREATE|os.O_RDWR, 0770)
}

func (d *dbLoc) updateFile(content []byte) error {
	d.Lock()
	defer d.Unlock()

	f, err := os.OpenFile(dbLocationFile+".tmp", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0770)
	if err != nil {
		return err
	}
	// cleanup our mess:
	// 1. Close might double-close the os.File, this should be fine in our context.
	// 2. Remove will try to remove the file even if renaming was successful, since we
	// only allow one caller at a time in here, it's safe for in-process consistency.
	defer os.Remove(f.Name())
	defer f.Close()

	err = f.Truncate(int64(len(content)))
	if err != nil {
		return err
	}

	n, err := f.Write(content)
	if err != nil {
		return err
	} else if n != len(content) {
		return errors.New("failed to write full content")
	}

	name := f.Name()
	if err = f.Close(); err != nil {
		return err
	}

	return os.Rename(name, dbLocationFile)
}

// Init initializes the config package, this should be called before any other
// functions are used in this package.
func Init() error {
	f, err := dbLocation.openFile()
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(&io.LimitedReader{R: f, N: maxLocLength})
	if err != nil {
		return err
	}

	if len(b) > 0 {
		dbLocation.path = string(b)
		return nil
	}

	return UpdateLocation(dbDefaultLocation)
}

// Location returns the current location of the database
func Location() string {
	dbLocation.Lock()
	defer dbLocation.Unlock()
	return dbLocation.path
}

// UpdateLocation updates the database location in the dbloc file, this
// does not move any existing databases around.
func UpdateLocation(path string) error {
	if len(path) > maxLocLength {
		return errors.New("new database location is too long")
	}

	err := dbLocation.updateFile([]byte(path))
	if err != nil {
		return err
	}

	dbLocation.path = path
	return nil
}
