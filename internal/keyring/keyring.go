package keyring

import (
	"fmt"

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

	// TEST TEST TEST TEST TEST TEST TEST TEST TEST TEST
	// I'll keep this code here until we get the delete function

	fmt.Printf("getting keyring: %s\n", name)
	i, err := getKeyring(ring, name)
	if err != nil {
		fmt.Println("err -", err)
	}
	fmt.Printf("keyring data: %s\n", i.Data)

	fmt.Printf("deleting keyring: %s\n", name)
	err = ring.Remove(name)
	if err != nil {
		fmt.Println("err -", err)
	}

	// TEST TEST TEST TEST TEST TEST TEST TEST TEST TEST

	return nil
}

func getKeyringService() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName: "papyro",
	})
}

func getKeyring(ring keyring.Keyring, name string) (keyring.Item, error) {
	return ring.Get(name)
}

func setKeyring(ring keyring.Keyring, name string, password string) error {
	return ring.Set(keyring.Item{
		Key:  name,
		Data: []byte(password),
	})
}
