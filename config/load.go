package config

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

const BucketName = "config"

// Config is a helper for configuration storage and retrieval.
//
// Configuration is stored in a boltdb bucket named after the
// BucketName constant. Storing and retrieval are handled by
// the various methods defined on this type.
//
// The current implementation of storing and retrieval uses
// encoding/json to encode and decode values passed to Load
// and Store.
type Config struct {
	name string
}

// OpenConfig returns the Config associated with the name passed in.
func OpenConfig(name string) Config {
	return Config{name}
}

// Load retrieves the value stored under the key and unmarshals it into
// the value given. If the key does not exist or is empty this returns
// err == nil.
func (c Config) Load(key string, value interface{}) error {
	db, err := Database(c.name)
	if err != nil {
		return err
	}

	return db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(BucketName))
		if buck == nil {
			return nil
		}

		b := buck.Get([]byte(key))
		// nothing to load, we don't treat this as an error
		if len(b) == 0 || b == nil {
			return nil
		}

		return json.Unmarshal(b, value)
	})
}

// Store stores the value given under the key passed. Store encodes the value
// given before storing it. The encoding used can be found in the Config
// documentation.
func (c Config) Store(key string, value interface{}) error {
	if value == nil {
		panic("config: Store value is nil")
	}

	// encode the value first, if it errors we don't need to touch the
	// database for writing.
	b, err := json.Marshal(value)
	if err != nil {
		return err
	} else if len(b) == 0 {
		// nothing to store in this case
		return nil
	}

	db, err := Database(c.name)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		buck, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			return err
		}

		return buck.Put([]byte(key), b)
	})
}
