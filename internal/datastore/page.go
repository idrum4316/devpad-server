package datastore

import (
	"encoding/json"
	"time"

	"github.com/idrum4316/devpad-server/internal/page"
	bolt "go.etcd.io/bbolt"
)

// UpdatePage updates a page in the datastore
func (d *Datastore) UpdatePage(p *page.Page, pageID string) error {

	p.Metadata.Modified = time.Now()

	pageBytes, err := json.Marshal(p)
	if err != nil {
		return err
	}

	err = d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pagesBucket))
		err := b.Put([]byte(pageID), pageBytes)
		return err
	})

	return err
}

// DeletePage deletes a page from the datastore.
func (d *Datastore) DeletePage(id string) error {

	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pagesBucket))
		err := b.Delete([]byte(id))
		return err
	})
	return err

}

// GetPage returns a page from the datastore
func (d *Datastore) GetPage(id string) (*page.Page, error) {

	var pageBytes []byte

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pagesBucket))
		v := b.Get([]byte(id))

		// This is not the most efficient way to do this
		if v != nil {
			pageBytes = append([]byte{}, v...)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if pageBytes == nil {
		return nil, nil
	}

	p := page.Page{}
	err = json.Unmarshal(pageBytes, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil

}
