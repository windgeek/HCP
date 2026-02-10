package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/spf13/cobra"
	"github.com/windgeek/HCP/pkg/config"
	"github.com/windgeek/HCP/pkg/identity"
	"github.com/windgeek/HCP/pkg/manifest"
)

var signCmd = &cobra.Command{
	Use:   "sign <file>",
	Short: "Sign a file with your identity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		// 1. Load Config
		keyPath, _ := cmd.Flags().GetString("key")
		cfg, err := config.LoadConfig(keyPath)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}
		identityPath := cfg.IdentityKeyPath

		if _, err := os.Stat(identityPath); os.IsNotExist(err) {
			fmt.Printf("Identity not found at %s. Please run 'hcp keygen' or check config.\n", identityPath)
			os.Exit(1)
		}

		fmt.Print("Enter passphrase: ")
		var passphrase string
		fmt.Scanln(&passphrase)

		key, err := identity.LoadKey(identityPath, passphrase)
		if err != nil {
			fmt.Printf("Error loading key: %v\n", err)
			os.Exit(1)
		}

		// 2. Get Author Address
		pubKey := key.PubKey()
		address, err := identity.PubKeyToAddress(pubKey, &chaincfg.MainNetParams)
		if err != nil {
			fmt.Printf("Error deriving address: %v\n", err)
			os.Exit(1)
		}
		
		pubKeyHex := hex.EncodeToString(pubKey.SerializeCompressed())

		// 3. Create Manifest
		m, err := manifest.NewManifest(filePath, address, pubKeyHex)
		if err != nil {
			fmt.Printf("Error creating manifest: %v\n", err)
			os.Exit(1)
		}

		// 4. Sign Manifest
		if err := m.Sign(key); err != nil {
			fmt.Printf("Error signing manifest: %v\n", err)
			os.Exit(1)
		}

		// 5. Save Manifest
		manifestPath := filePath + ".hcp"
		if err := m.Save(manifestPath); err != nil {
			fmt.Printf("Error saving manifest: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("File signed successfully: %s\n", manifestPath)
	},
}

func init() {
	signCmd.Flags().String("key", "", "Path to identity key file")
	rootCmd.AddCommand(signCmd)
}
