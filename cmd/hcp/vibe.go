package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/windgeek/HCP/pkg/entropy"
	"golang.org/x/term"
)

var vibeCmd = &cobra.Command{
	Use:   "vibe",
	Short: "Capture biological entropy (Proof of Hesitation)",
	Long:  `Interactive session to capture keystroke dynamics and generate a Biological Entropy Proof.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("HCP Biological Entropy Engine (BEE)")
		fmt.Println("-----------------------------------")
		fmt.Println("Please type the following phrase exactly as shown, then press Enter:")
		targetPhrase := "The quick brown fox jumps over the lazy dog."
		fmt.Printf("\n> %s\n\n", targetPhrase)
		fmt.Print("Start typing: ")

		// Set terminal to raw mode
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Printf("Error setting raw mode: %v\n", err)
			os.Exit(1)
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)

		var keystrokes []entropy.Keystroke
		var inputString string

		// Capture Loop
		buf := make([]byte, 1)
		for {
			_, err := os.Stdin.Read(buf)
			if err != nil {
				break
			}
			char := rune(buf[0])
			now := time.Now()

			// Handle Ctrl+C (3) or Enter (13)
			if char == 3 {
				term.Restore(int(os.Stdin.Fd()), oldState)
				os.Exit(1)
			}
			if char == 13 || char == 10 { // Enter
				break
			}

			// Handle Backspace (127)
			if char == 127 {
				if len(inputString) > 0 {
					inputString = inputString[:len(inputString)-1]
					// Remove last keystroke? RFC-002 says "Refactoring Entropy" is valuable.
					// We should KEEP the backspace keystroke in the log to show "modification".
					// But for visual feedback we erase.
					fmt.Print("\b \b")
				}
			} else {
				inputString += string(char)
				fmt.Print(string(char))
			}

			keystrokes = append(keystrokes, entropy.Keystroke{
				Key:       char,
				Timestamp: now,
			})
		}

		// Restore terminal
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Println("\n\nAnalyzing entropy...")

		// Analyze
		stats, err := entropy.Analyze(keystrokes)
		if err != nil {
			fmt.Printf("Error analyzing entropy: %v\n", err)
			os.Exit(1)
		}

		// Generate Hash
		sessionHash := entropy.GenerateSessionHash(keystrokes)

		// Output Results
		fmt.Println("-----------------------------------")
		fmt.Printf("Total Keystrokes: %d\n", stats.TotalKeystrokes)
		fmt.Printf("Duration:         %.2fs\n", stats.Duration)
		fmt.Printf("Mean Flight Time: %.2fms\n", stats.MeanFlightTime)
		fmt.Printf("Variance:         %.2f\n", stats.Variance)
		fmt.Printf("Shannon Entropy:  %.2f bits\n", stats.ShannonEntropy)
		fmt.Printf("Human Score:      %.2f / 1.0\n", stats.HumanScore)
		fmt.Println("-----------------------------------")
		fmt.Printf("Session Manifest Hash: %s\n", sessionHash)

		// Save Session Data
		sessionData := struct {
			Target  string                `json:"target_phrase"`
			Input   string                `json:"input_phrase"`
			Stats   *entropy.SessionStats `json:"stats"`
			RawHash string                `json:"session_hash"`
		}{
			Target:  targetPhrase,
			Input:   inputString,
			Stats:   stats,
			RawHash: sessionHash,
		}

		jsonData, _ := json.MarshalIndent(sessionData, "", "  ")
		if err := os.WriteFile("hcp-session.json", jsonData, 0644); err == nil {
			fmt.Println("Session data saved to hcp-session.json")
		}
	},
}

func init() {
	rootCmd.AddCommand(vibeCmd)
}
