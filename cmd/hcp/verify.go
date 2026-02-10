package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/spf13/cobra"
	"github.com/windgeek/HCP/pkg/identity"
	"github.com/windgeek/HCP/pkg/manifest"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify the integrity and authorship of the current directory",
	Long:  `Verify that the current directory matches the manifest.hcp and that the signature is valid.`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		manifestPath := filepath.Join(cwd, "manifest.hcp")
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			fmt.Println("Error: manifest.hcp not found in current directory.")
			os.Exit(1)
		}

		// 1. Read Manifest
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			fmt.Printf("Error reading manifest: %v\n", err)
			os.Exit(1)
		}

		var m manifest.Manifest
		if err := json.Unmarshal(data, &m); err != nil {
			fmt.Printf("Error parsing manifest: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Verifying Manifest...")
		fmt.Printf("Author: %s\n", m.Author)
		fmt.Printf("Timestamp: %d\n", m.Timestamp)

		// 2. Verify PublicKey matches Author Address
		pubKeyBytes, err := hex.DecodeString(m.PublicKey)
		if err != nil {
			fmt.Printf("Error decoding public key: %v\n", err)
			os.Exit(1)
		}
		pubKey, err := btcec.ParsePubKey(pubKeyBytes)
		if err != nil {
			fmt.Printf("Error parsing public key: %v\n", err)
			os.Exit(1)
		}

		currAddr, err := identity.PubKeyToAddress(pubKey, &chaincfg.MainNetParams)
		if err != nil {
			fmt.Printf("Error deriving address: %v\n", err)
			os.Exit(1)
		}

		if currAddr != m.Author {
			fmt.Printf("[FAIL] Public Key does not match Author Address!\nDetailed: Derived %s vs Claimed %s\n", currAddr, m.Author)
			os.Exit(1)
		}
		fmt.Println("[PASS] Author Identity Verified")

		// 3. Verify Signature
		// Reconstruct payload hash
		// We need to replicate the exact payload struct from pkg/manifest/manifest.go
		// Since we can't easily access the private payload struct or Sign logic without key,
		// we should probably expose a Verify method in manifest package.
		// For now, let's implement signature verification here by duplicating the struct logic
		// OR better: Refactor Verify into pkg/manifest.

		// Let's rely on pkg/manifest having a Verify method. 
		// I'll add Verify to pkg/manifest in next step. For now, assuming it exists or implementing logic here.
		
		if err := verifySignature(&m, pubKey); err != nil {
			fmt.Printf("[FAIL] Invalid Signature: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("[PASS] Cryptographic Signature Verified")

		// 4. Verify Content Integrity
		fmt.Println("Verifying Content Integrity...")
		// Load ignores
		ignorePatterns := []string{".git", ".hcp", "node_modules", ".DS_Store", "*.hcp"}
		// TODO: Load .hcpignore if exists
		
		calcHash, _, _, _, err := manifest.CalculateDirHash(cwd, ignorePatterns)
		if err != nil {
			fmt.Printf("Error calculating hash: %v\n", err)
			os.Exit(1)
		}

		if calcHash != m.ContentHash {
			fmt.Printf("[WARNING] Fingerprint Mismatch!\n")
			fmt.Printf("Manifest Hash:   %s\n", m.ContentHash)
			fmt.Printf("Calculated Hash: %s\n", calcHash)
			fmt.Println("This file has been altered by non-sovereign entities.")
			os.Exit(1)
		}

		fmt.Println("[SUCCESS] Human Intent Verified. Integrity 100%.")
	},
}

func verifySignature(m *manifest.Manifest, pubKey *btcec.PublicKey) error {
	// Reconstruct payload
	// NOTE: This must exact match pkg/manifest/manifest.go's payload
	// It's better to move this verification logic to `pkg/manifest`.
	// I will implement Verify in pkg/manifest/manifest.go and call it.
	return m.Verify(pubKey)
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
