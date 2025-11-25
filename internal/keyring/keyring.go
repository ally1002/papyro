package keyring

import (
	"github.com/99designs/keyring"
)

func SavePassword(name string, password string) error {
	ring, err := getKeyringService()
	if err != nil {
		return err
	}

	err = setKeyring(ring, name, password)
	if err != nil {
		return err
	}

	return nil
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
