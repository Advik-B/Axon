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
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Add the 'build' command to the root command.
	rootCmd.AddCommand(buildCmd)
}