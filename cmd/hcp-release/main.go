package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/windgeek/HCP/pkg/identity"
	"github.com/windgeek/HCP/pkg/manifest"
	"golang.org/x/term"
)

func main() {
	// 1. Parse Flags
	targetPath := flag.String("path", ".", "Path to the directory to release")
	flag.Parse()

	// Resolve absolute path for scanning
	absPath, err := filepath.Abs(*targetPath)
	if err != nil {
		fmt.Printf("Error resolving path: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("HCP Release Manifest Generator")
	fmt.Println("------------------------------")
	// Calculate relative path for display
	cwd, _ := os.Getwd()
	displayTarget := absPath
	if rel, err := filepath.Rel(cwd, absPath); err == nil {
		displayTarget = rel
	}
	fmt.Printf("Target Path: %s\n", displayTarget)

	// 2. Load .hcpignore
	ignorePatterns := loadIgnorePatterns(absPath)
	// Add default ignores
	ignorePatterns = append(ignorePatterns, ".git", ".hcp", "node_modules", ".DS_Store", "*.hcp")

	// 3. Scan and Hash
	globalHash, assets, contribMap, zkpMap, err := manifest.CalculateDirHash(absPath, ignorePatterns)
	if err != nil {
		fmt.Printf("Error calculating hash: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Total Assets: %d\n", len(assets))
	
	var totalScore float64
	for _, m := range contribMap {
		totalScore += m.AHAScore
	}
	avgScore := 0.0
	if len(assets) > 0 {
		avgScore = totalScore / float64(len(assets))
	}
	
	fmt.Printf("Global Content Hash: %s\n", globalHash)
	fmt.Printf("Average AHA Score: %.1f / 100\n", avgScore)

	// 4. Determine Output Filename
	defaultFilename := "manifest.hcp"
	defaultOutputPath := filepath.Join(absPath, defaultFilename)
	finalOutputPath := defaultOutputPath

	if _, err := os.Stat(defaultOutputPath); err == nil {
		// File exists, prompt for version
		
		// Check if interactive
		stat, _ := os.Stdin.Stat()
		isInteractive := (stat.Mode() & os.ModeCharDevice) != 0

		if isInteractive {
			fmt.Printf("\n'%s' already exists.\n", defaultFilename)
			fmt.Print("Enter version tag to create a new file (e.g., 'v1.0'), or press Enter to overwrite: ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input != "" {
				// Sanitize input
				tag := strings.ReplaceAll(input, " ", "_")
				finalOutputPath = filepath.Join(absPath, fmt.Sprintf("manifest-%s.hcp", tag))
			}
		} 
		// If piped, we default to overwrite for automation, unless logic dictates otherwise. 
		// For now, automation overwrites.
	}

	// 5. Load Identity
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

	var passphrase string
	fmt.Print("\nEnter passphrase to sign release: ")
	
	// Check if simple stdin (piped) or terminal
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Piped
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			passphrase = scanner.Text()
		}
		fmt.Println()
	} else {
		// Interactive
		passBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Printf("\nError reading password: %v\n", err)
			os.Exit(1)
		}
		passphrase = string(passBytes)
		fmt.Println()
	}

	key, err := identity.LoadKey(identityPath, passphrase)
	if err != nil {
		fmt.Printf("Error loading key: %v\n", err)
		os.Exit(1)
	}

	authAddr, err := identity.PubKeyToAddress(key.PubKey(), &chaincfg.MainNetParams)
	if err != nil {
		fmt.Printf("Error deriving address: %v\n", err)
		os.Exit(1)
	}
	
	pubKeyHex := hex.EncodeToString(key.PubKey().SerializeCompressed())

	// 6. Create Manifest
	m := manifest.Manifest{
		Version:     "v1-release",
		Author:      authAddr,
		PublicKey:   pubKeyHex,
		ContentHash: globalHash,
		Timestamp:   time.Now().Unix(),
		EntropyDNA:      "universal-release",
		Assets:          assets,
		ContributionMap: contribMap,
		CognitiveProofs: zkpMap,
	}

	// 7. Sign
	if err := m.Sign(key); err != nil {
		fmt.Printf("Error signing manifest: %v\n", err)
		os.Exit(1)
	}

	// 8. Save
	if err := m.Save(finalOutputPath); err != nil {
		fmt.Printf("Error saving manifest: %v\n", err)
		os.Exit(1)
	}
	// 9. Format Output Path for Display
	cwd, _ = os.Getwd()
	displayPath := finalOutputPath
	if rel, err := filepath.Rel(cwd, finalOutputPath); err == nil {
		if !strings.HasPrefix(rel, "..") {
			displayPath = rel
		}
	}
	fmt.Printf("\nRelease Manifest generated: %s\n", displayPath)
}

func loadIgnorePatterns(root string) []string {
	var patterns []string
	f, err := os.Open(filepath.Join(root, ".hcpignore"))
	if err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				patterns = append(patterns, line)
			}
		}
	}
	return patterns
}

