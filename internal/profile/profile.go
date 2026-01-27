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
	config   *config.Config
}

func NewProfiles() (*Profiles, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	ps := &Profiles{config: cfg}

	if err := ps.load(); err != nil {
		ps.Profiles = []Profile{}
	}

	return ps, nil
}

func (ps *Profiles) Get(name string) (*Profile, error) {
	for _, profile := range ps.Profiles {
		if profile.Name == name {
			return &profile, nil
		}
	}

	return &Profile{}, fmt.Errorf("profile '%s' does not exist", name)
}

func (ps *Profiles) Add(p *Profile) error {
	if err := p.Validate(); err != nil {
		return err
	}

	p, err := ps.Get(p.Name)
	if err == nil {
		return fmt.Errorf("profile '%s' already exists", p.Name)
	}

	ps.Profiles = append(ps.Profiles, *p)

	return ps.save()
}

func ReadProfiles() error {
	data, err := getProfiles()
	if err != nil {
		return err
	}

	return profilesTable(data.Profiles)
}

func (ps *Profiles) Delete(name string) error {
	profile, err := ps.Get(name)
	if err != nil {
		return err
	}

	profileIndex := slices.Index(ps.Profiles, *profile)

	ps.Profiles = slices.Delete(ps.Profiles, profileIndex, profileIndex+1)

	profiles, err := json.MarshalIndent(ps.Profiles, "", "  ")
	if err != nil {
		return err
	}

	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("could not determine user config directory: %w", err)
	}

	return os.WriteFile(cfg.FilePath, profiles, 0600)
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

func (ps *Profiles) save() error {
	data, err := json.MarshalIndent(ps.Profiles, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ps.config.FilePath, data, 0600)
}

func (ps *Profiles) load() error {
	file, err := os.ReadFile(ps.config.FilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &ps.Profiles)
	if err != nil {
		return err
	}

	return nil
}

func (p *Profile) Validate() error {
	if p.Name == "" && p.FromEmail == "" && p.KindleEmail == "" {
		return fmt.Errorf("required args cannot be blank")
	}

	return nil
}
