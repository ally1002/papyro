package keyring

import (
	"github.com/99designs/keyring"
)

type Keyring struct {
	ring keyring.Keyring
}

func NewRing() (*Keyring, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "papyro",
	})
	if err != nil {
		return &Keyring{}, err
	}

	return &Keyring{ring: ring}, nil
}

func (kr *Keyring) Save(name string, password string) error {
	return setKeyring(kr.ring, name, password)
}

func DeletePassword(name string) error {
	ring, err := getKeyringService()
	if err != nil {
		return err
	}

	err = ring.Remove(name)
	if err != nil {
		return err
	}

	return nil
}

func getKeyringService() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName: "papyro",
	})
}

func setKeyring(ring keyring.Keyring, name string, password string) error {
	return ring.Set(keyring.Item{
		Key:  name,
		Data: []byte(password),
	})
}
