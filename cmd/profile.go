/*
Copyright © 2025 ally1002
*/
package cmd

import (
	"fmt"

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
		fmt.Println(password)

		return profile.WriteProfile(name, fromEmail, kindleEmail)
	},
}

func init() {
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
	if err := profileAddCmd.MarkFlagRequired("kindle-email"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileAddCmd)
}
