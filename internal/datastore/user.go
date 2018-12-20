package datastore

import (
	"encoding/json"
	"errors"

	"github.com/idrum4316/devpad-server/internal/user"
	bolt "go.etcd.io/bbolt"
)

// CountUsers returns the number of users in the datastore
func (d *Datastore) CountUsers() (int, error) {
	count := 0

	err := d.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(usersBucket))

		// There's no code that could produce an error to return
		_ = b.ForEach(func(k, v []byte) error {
			count = count + 1
			return nil
		})

		return nil
	})

	return count, err
}

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

// CreateUser creates a new user in the database
func (d *Datastore) CreateUser(u *user.User) error {

	exists, err := d.UserExists(u.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}

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

// UpdateUser updates a user in the datastore
func (d *Datastore) UpdateUser(u *user.User) error {

	exists, err := d.UserExists(u.ID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user does not exist")
	}

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
