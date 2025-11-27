package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"text/tabwriter"

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

func CreateProfile(profile *Profile) error {
	if err := isProfileValid(profile); err != nil {
		return err
	}

	profiles, err := writeProfile(profile)
	if err != nil {
		return err
	}

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	err = os.WriteFile(cfg.FilePath, profiles, 0600)

	return err
}

func ReadProfiles() error {
	data, err := getProfiles()
	if err != nil {
		return err
	}

	return profilesTable(data.Profiles)
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

	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("could not determine user config directory: %w", err)
	}

	return os.WriteFile(cfg.FilePath, profiles, 0600)
}

func isProfileValid(profile *Profile) error {
	if profile.Name == "" && profile.FromEmail == "" && profile.KindleEmail == "" {
		return fmt.Errorf("required args cannot be blank")
	}

	_, err := getProfile(profile.Name)
	if err == nil {
		return fmt.Errorf("profile '%s' already exists", profile.Name)
	}

	return nil
}

func getProfiles() (Profiles, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return Profiles{}, fmt.Errorf("could not determine user config directory: %w", err)
	}

	file, err := os.ReadFile(cfg.FilePath)
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

func writeProfile(profile *Profile) ([]byte, error) {
	data, err := getProfiles()
	if err != nil {
		return nil, err
	}

	data.Profiles = append(data.Profiles, *profile)

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

func profilesTable(profiles []Profile) error {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	_, err := fmt.Fprintln(writer, "NAME\tFROM EMAIL\tKINDLE EMAIL")
	if err != nil {
		return err
	}

	for _, p := range profiles {
		_, err := fmt.Fprintf(writer, "%s\t%s\t%s\n", p.Name, p.FromEmail, p.KindleEmail)
		if err != nil {
			return err
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
