package config

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

const BucketName = "config"

// Config is a helper for configuration storage and retrieval.
type Config struct {
	name string
}

// OpenConfig returns the Config associated with the name given.
func OpenConfig(name string) Config {
	return Config{name}
}

// Load retrieves the value stored under the key and unmarshals it into
// the value given. If the key does not exist or is empty this returns
// err == nil.
//
// Load uses an encoding to load the value, any fields in the
// value type need to be exported to load properly.
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
//
// Store uses an encoding to store the value given, this means that any fields
// in the value type need to be exported to be saved.
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
