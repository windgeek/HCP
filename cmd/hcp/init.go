package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/windgeek/HCP/pkg/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize HCP configuration",
	Long:  `Generate a default .hcp/config.yaml file in the current directory or home directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		global, _ := cmd.Flags().GetBool("global")
		
		var configPath string
		if global {
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Printf("Error getting home directory: %v\n", err)
				os.Exit(1)
			}
			configPath = filepath.Join(home, ".hcp", "config.yaml")
		} else {
			configPath = filepath.Join(".hcp", "config.yaml")
		}

		cfg, err := config.DefaultConfig()
		if err != nil {
			fmt.Printf("Error creating default config: %v\n", err)
			os.Exit(1)
		}

		if err := config.SaveConfig(configPath, cfg); err != nil {
			fmt.Printf("Error saving config to %s: %v\n", configPath, err)
			os.Exit(1)
		}

		fmt.Printf("Initialized HCP configuration at: %s\n", configPath)
		fmt.Printf("Identity Key Path: %s\n", cfg.IdentityKeyPath)
	},
}

func init() {
	initCmd.Flags().BoolP("global", "g", false, "Initialize in user home directory (~/.hcp/config.yaml)")
	rootCmd.AddCommand(initCmd)
}
