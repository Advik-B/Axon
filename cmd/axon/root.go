package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "axon",
	Short: "Axon is a visual, node-based programming language that transpiles to Go.",
	Long: `A fast and flexible tool to build, run, and manage Axon visual programs.
Axon transpiles .ax graph files into readable, idiomatic Go code.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Add the subcommands to the root command.
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(packCmd)
	rootCmd.AddCommand(unpackCmd)
	rootCmd.AddCommand(convertCmd)
}