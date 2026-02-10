package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/windgeek/HCP/pkg/manifest"
)

var anchorCmd = &cobra.Command{
	Use:   "anchor <file.hcp>",
	Short: "Anchor a signed manifest to Bitcoin (Mock)",
	Long:  `Simulates anchoring a signed manifest to the Bitcoin blockchain using OP_RETURN.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manifestPath := args[0]

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

		// 2. Generate OP_RETURN Data
		// OP_RETURN (0x6a) + Length (0x20 = 32 bytes) + ContentHash (32 bytes)
		// ContentHash is already hex string in manifest, we need to decode it to verify length/bytes if we were constructing raw script.
		// For display, we just concatenate.
		// Check hash length
		hashBytes, err := hex.DecodeString(m.ContentHash)
		if err != nil || len(hashBytes) != 32 {
			fmt.Printf("Invalid content hash in manifest: %s\n", m.ContentHash)
			os.Exit(1)
		}

		opReturnHex := fmt.Sprintf("6a20%s", m.ContentHash)

		// 3. Log to file
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		logPath := filepath.Join(home, ".hcp", "anchor.log")

		logEntry := fmt.Sprintf("[%s] Anchored %s: %s\n", time.Now().Format(time.RFC3339), manifestPath, opReturnHex)

		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening log file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()

		if _, err := f.WriteString(logEntry); err != nil {
			fmt.Printf("Error writing to log: %v\n", err)
			os.Exit(1)
		}

		// 4. Output to User
		fmt.Printf("Mock anchoring to Bitcoin: %s\n", opReturnHex)
		fmt.Printf("Transaction logged to %s\n", logPath)
	},
}

func init() {
	rootCmd.AddCommand(anchorCmd)
}
