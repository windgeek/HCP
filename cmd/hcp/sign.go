package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/spf13/cobra"
	"github.com/windgeek/HCP/pkg/identity"
	"github.com/windgeek/HCP/pkg/manifest"
)

var signCmd = &cobra.Command{
	Use:   "sign <file>",
	Short: "Sign a file with your identity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		// 1. Load Identity
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		identityPath := filepath.Join(home, ".hcp", "identity.key")

		if _, err := os.Stat(identityPath); os.IsNotExist(err) {
			fmt.Println("Identity not found. Please run 'hcp keygen' first.")
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
	rootCmd.AddCommand(signCmd)
}
