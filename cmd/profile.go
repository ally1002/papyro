/*
Copyright © 2025 ally1002
*/
package cmd

import (
	"fmt"

	"github.com/ally1002/papyro/internal/keyring"
	"github.com/ally1002/papyro/internal/profile"
	"github.com/spf13/cobra"
)

var name, kindleEmail, fromEmail, password string

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage profiles",
}

var profileAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		ps, err := profile.NewProfiles()
		if err != nil {
			return err
		}

		kr, err := keyring.NewRing()
		if err != nil {
			return err
		}

		err = ps.Add(&profile.Profile{Name: name, FromEmail: fromEmail, KindleEmail: kindleEmail})
		if err != nil {
			return err
		}

		err = kr.Save(name, password)
		if err != nil {
			return err
		}

		return nil
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		ps, err := profile.NewProfiles()
		if err != nil {
			return err
		}

		return ps.List()
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf(`required argument "name" was not set`)
		}

		profileName := args[0]

		ps, err := profile.NewProfiles()
		if err != nil {
			return err
		}

		kr, err := keyring.NewRing()
		if err != nil {
			return err
		}

		err = ps.Delete(profileName)
		if err != nil {
			return err
		}

		err = kr.Delete(profileName)
		if err != nil {
			return err
		}

		// need to add a here rollback later

		return nil
	},
}

func init() {
	profileAddFlags()

	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileDeleteCmd)
}

func profileAddFlags() {
	profileAddCmd.Flags().StringVarP(&name, "name", "n", "", "user name")
	if err := profileAddCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	profileAddCmd.Flags().StringVarP(&kindleEmail, "kindle-email", "k", "", "kindle email")
	if err := profileAddCmd.MarkFlagRequired("kindle-email"); err != nil {
		panic(err)
	}

	profileAddCmd.Flags().StringVarP(&fromEmail, "from-email", "f", "", "sender email")
	if err := profileAddCmd.MarkFlagRequired("from-email"); err != nil {
		panic(err)
	}

	profileAddCmd.Flags().StringVar(&password, "passwd", "", "sender email password")
	if err := profileAddCmd.MarkFlagRequired("passwd"); err != nil {
		panic(err)
	}
}
