/*
Copyright © 2025 ally1002
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ally1002/papyro/internal/email"
	"github.com/ally1002/papyro/internal/keyring"
	"github.com/ally1002/papyro/internal/profile"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send attachments to kindle",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("2 required arguments: profileName, filePath")
		}

		profileName := args[0]
		filePath := args[1]
		if !filepath.IsAbs(filePath) {
			pwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}
			filePath = filepath.Join(pwd, filePath)
		}

		ps, err := profile.NewProfiles()
		if err != nil {
			return err
		}

		prof, err := ps.Get(profileName)
		if err != nil {
			return err
		}

		kr, err := keyring.NewRing()
		if err != nil {
			return err
		}

		item, err := kr.Get(profileName)
		if err != nil {
			return err
		}

		password := item.Data

		mail, err := email.NewEmail(*prof, password, filePath)
		if err != nil {
			return err
		}

		if err := mail.Send(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
