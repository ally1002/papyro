package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/ally1002/papyro/internal/config"
)

type Profile struct {
	Name        string `json:"name"`
	FromEmail   string `json:"fromEmail"`
	KindleEmail string `json:"kindleEmail"`
}

type Profiles struct {
	Profiles []Profile `json:"profiles"`
}

func CreateProfile(name string, fromEmail string, kindleEmail string) error {
	if err := isProfileValid(name, fromEmail, kindleEmail); err != nil {
		return err
	}

	profiles, err := writeProfile(name, fromEmail, kindleEmail)
	if err != nil {
		return err
	}

	_, fileName, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("could not determine user home directory: %w", err)
	}

	err = os.WriteFile(fileName, profiles, 0600)

	return err
}

func DeleteProfile(name string) error {
	profile, err := getProfile(name)
	if err != nil {
		return err
	}

	data, err := getProfiles()
	if err != nil {
		return err
	}

	profileIndex := slices.Index(data.Profiles, profile)
	data.Profiles = slices.Delete(data.Profiles, profileIndex, profileIndex+1)

	profiles, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	_, fileName, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("could not determine user home directory: %w", err)
	}

	return os.WriteFile(fileName, profiles, 0600)
}

func isProfileValid(name string, fromEmail string, kindleEmail string) error {
	if name == "" && fromEmail == "" && kindleEmail == "" {
		return fmt.Errorf("required args cannot be blank")
	}

	_, err := getProfile(name)
	if err == nil {
		return fmt.Errorf("profile '%s' already exists", name)
	}

	return nil
}

func getProfiles() (Profiles, error) {
	_, fileName, err := config.GetConfig()
	if err != nil {
		return Profiles{}, fmt.Errorf("could not determine user home directory: %w", err)
	}

	file, err := os.ReadFile(fileName)
	if err != nil {
		return Profiles{}, err
	}

	var data Profiles
	err = json.Unmarshal(file, &data)
	if err != nil {
		return Profiles{}, err
	}

	return data, nil
}

func writeProfile(name string, fromEmail string, kindleEmail string) ([]byte, error) {
	data, err := getProfiles()
	if err != nil {
		return nil, err
	}

	data.Profiles = append(data.Profiles, Profile{Name: name, FromEmail: fromEmail, KindleEmail: kindleEmail})

	profiles, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func getProfile(name string) (Profile, error) {
	data, err := getProfiles()
	if err != nil {
		return Profile{}, err
	}

	for _, profile := range data.Profiles {
		if profile.Name == name {
			return profile, nil
		}
	}

	return Profile{}, fmt.Errorf("profile '%s' does not exist", name)
}
