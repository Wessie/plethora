package plethora

import (
	"time"

	"github.com/Wessie/plethora/config"
	"github.com/boltdb/bolt"
)

const linkEntryBucket = "links_by_entry"
const linkSourceBucket = "links_by_source"

// StoreData is equal to StoreIdentifier but removes a few steps.
func StoreData(d Data, sourceTime time.Time) error {
	return StoreIdentifier(Identifier(d), sourceTime)
}

// StoreIdentifier stores the identifier, it should point to valid data.
// The sourceTime is the time this data was created at the data source
// (e.g news article posting date/time).
func StoreIdentifier(id identifier, sourceTime time.Time) error {
	db, err := config.Database("plethora")
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(linkEntryBucket))
		if err != nil {
			return err
		}

		entryTime := time.Now().Format(time.RFC3339Nano)
		err = b.Put([]byte(entryTime), []byte(id))
		if err != nil {
			return err
		}

		b, err = tx.CreateBucketIfNotExists([]byte(linkSourceBucket))
		if err != nil {
			return err
		}

		srcTime := sourceTime.Format(time.RFC3339Nano)
		return b.Put([]byte(srcTime), []byte(id))
	})
}
