package user

import (
	"bytes"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/scrypt"
)

const (
	pwSaltBytes = 32
	pwHashBytes = 64
)

// User is an API user - or web user
type User struct {
	ID       string
	Password []byte
	Salt     []byte
	Admin    bool
}

// SetPassword generates a new salt and uses it to has the user's password
func (u *User) SetPassword(password string) error {

	salt := make([]byte, pwSaltBytes)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return err
	}

	var derivedKey []byte
	derivedKey, err = scrypt.Key([]byte(password), salt, 1<<15, 8, 1, pwHashBytes)
	if err != nil {
		return err
	}

	u.Password = derivedKey
	u.Salt = salt

	return nil

}

// VerifyPassword checks to see if the password parameter matches the User's
// password
func (u *User) VerifyPassword(password string) (bool, error) {

	derivedKey, err := scrypt.Key([]byte(password), u.Salt, 1<<15, 8, 1, pwHashBytes)
	if err != nil {
		return false, err
	}

	return bytes.Equal(derivedKey, u.Password), nil

}
