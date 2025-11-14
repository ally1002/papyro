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
		err := profile.CreateProfile(name, fromEmail, kindleEmail)
		if err != nil {
			return err
		}

		err = keyring.SavePassword(name, password)
		if err != nil {
			return err
		}

		// need to add a rollback here after destroy is created

		return nil
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf(`required argument "name" was not set`)
		}

		err := profile.DeleteProfile(args[0])
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	profileAddFlags()

	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileAddCmd)
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
