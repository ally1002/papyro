/*
Copyright © 2025 ally1002
*/

package cmd

import (
	"log"
	"os"

	"github.com/ally1002/papyro/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "papyro",
	Short: "A multi-profile kindle sender",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	cfg.CheckAndCreateConfiguration()
}
