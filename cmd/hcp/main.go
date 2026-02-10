package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hcp",
	Short: "HCP: Human Core Protocol CLI",
	Long:  `HCP is a decentralized protocol designed to provide an immutable "Proof of Humanity" for digital assets.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
