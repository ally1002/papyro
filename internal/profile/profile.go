package profile

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
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
	idx := slices.IndexFunc(ps.Profiles, func(p Profile) bool {
		return p.Name == name
	})

	if idx == -1 {
		return nil, fmt.Errorf("profile '%s' does not exist", name)
	}

	return &ps.Profiles[idx], nil
}

func (ps *Profiles) Add(p *Profile) error {
	if err := p.Validate(); err != nil {
		return err
	}

	profile, err := ps.Get(p.Name)
	if err == nil {
		return fmt.Errorf("profile '%s' already exists", profile.Name)
	}

	ps.Profiles = append(ps.Profiles, *p)

	return ps.save()
}

func (ps *Profiles) List(w io.Writer) error {
	writer := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)

	_, err := fmt.Fprintln(writer, "NAME\tFROM EMAIL\tKINDLE EMAIL")
	if err != nil {
		return err
	}

	for _, p := range ps.Profiles {
		_, err := fmt.Fprintf(writer, "%s\t%s\t%s\n", p.Name, p.FromEmail, p.KindleEmail)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func (ps *Profiles) Delete(name string) error {
	profile, err := ps.Get(name)
	if err != nil {
		return err
	}

	profileIndex := slices.Index(ps.Profiles, *profile)

	ps.Profiles = slices.Delete(ps.Profiles, profileIndex, profileIndex+1)

	return ps.save()
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
	if p.Name == "" || p.FromEmail == "" || p.KindleEmail == "" {
		return fmt.Errorf("required args cannot be blank")
	}

	if !isValidEmail(p.FromEmail) {
		return fmt.Errorf("fromEmail is invalid")
	}

	if !isValidEmail(p.KindleEmail) {
		return fmt.Errorf("kindleEmail is invalid")
	}

	return nil
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
