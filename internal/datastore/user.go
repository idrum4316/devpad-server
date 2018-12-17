package datastore

import (
	"encoding/json"

	"github.com/idrum4316/devpad-server/internal/user"
	bolt "go.etcd.io/bbolt"
)

// UserExists returns a false if the user doesn't exist in the datastore
func (d *Datastore) UserExists(id string) (bool, error) {

	exists := false

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(usersBucket))
		v := b.Get([]byte(id))

		if v != nil {
			exists = true
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return exists, nil

}

// UpdateUser updates a user in the datastore
func (d *Datastore) UpdateUser(u *user.User) error {

	userBytes, err := json.Marshal(u)
	if err != nil {
		return err
	}

	err = d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(usersBucket))
		err := b.Put([]byte(u.ID), userBytes)
		return err
	})

	return err
}

// DeleteUser deletes a user from the datastore.
func (d *Datastore) DeleteUser(id string) error {

	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(usersBucket))
		err := b.Delete([]byte(id))
		return err
	})
	return err

}

// GetUser returns a page from the datastore
func (d *Datastore) GetUser(id string) (*user.User, error) {

	var userBytes []byte

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(usersBucket))
		v := b.Get([]byte(id))

		// This is not the most efficient way to do this
		if v != nil {
			userBytes = append([]byte{}, v...)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if userBytes == nil {
		return nil, nil
	}

	u := user.User{}
	err = json.Unmarshal(userBytes, &u)
	if err != nil {
		return nil, err
	}

	return &u, nil

}
