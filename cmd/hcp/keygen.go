package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/spf13/cobra"
	"github.com/windgeek/HCP/pkg/identity"
)

var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate a new identity key",
	Long:  `Generate a new secp256k1 private key and save it loosely encrypted to ~/.hcp/identity.key`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Determine path
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		identityPath := filepath.Join(home, ".hcp", "identity.key")

		// 2. Check if key already exists
		if _, err := os.Stat(identityPath); err == nil {
			fmt.Printf("Identity key already exists at %s\n", identityPath)
			// For safety, don't overwrite unless force flag (not implemented yet) is used.
			// But for this task, let's just warn and exit or prompt.
			// The prompt implies "Implement hcp keygen command".
			// I'll just return to accidental overwrites.
			return
		}

		// 3. Generate Key
		fmt.Println("Generating new identity key...")
		privKey, err := identity.GenerateKey()
		if err != nil {
			fmt.Printf("Error generating key: %v\n", err)
			os.Exit(1)
		}

		// 4. Get Passphrase (simple prompt)
		// For prototype, we can use a default or ask.
		// Let's ask.
		fmt.Print("Enter passphrase to encrypt your key: ")
		var passphrase string
		fmt.Scanln(&passphrase)
		if passphrase == "" {
			fmt.Println("Passphrase cannot be empty.")
			os.Exit(1)
		}

		// 5. Save Key
		if err := identity.SaveKey(privKey, identityPath, passphrase); err != nil {
			fmt.Printf("Error saving key: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Key saved to %s\n", identityPath)

		// 6. Output Address
		// Using MainNet for the address generation as default
		address, err := identity.PubKeyToAddress(privKey.PubKey(), &chaincfg.MainNetParams)
		if err != nil {
			fmt.Printf("Error deriving address: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Your Identity Address (P2WPKH): %s\n", address)
	},
}

func init() {
	rootCmd.AddCommand(keygenCmd)
}
