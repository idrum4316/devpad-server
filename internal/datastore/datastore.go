package datastore

import (
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	pagesBucket = "Pages"
	usersBucket = "Users"
)

// Datastore is where user accounts and page metadata is stored
type Datastore struct {
	db *bolt.DB
}

// New returns a new, already opened Datastore instance at <path>
func New(path string) (*Datastore, error) {

	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return nil, err
	}

	store := Datastore{
		db: db,
	}
	store.initialize()

	return &store, nil

}

// Initialize the buckets that are needed for future transactions
func (d *Datastore) initialize() (err error) {

	err = d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(usersBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(pagesBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return

}

// Close closes the enclosed bolt database
func (d *Datastore) Close() {
	d.db.Close()
}
