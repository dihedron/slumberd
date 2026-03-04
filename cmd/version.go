/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/dihedron/slumberd/metadata"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the command version and exit.",
	Long:  `Show the command version and exit.`,
	Run:   RunVersion,
}

func RunVersion(cmd *cobra.Command, args []string) {
	slog.Info("running version command")
	verbose, _ := cmd.Flags().GetBool("verbose")
	if verbose {
		metadata.PrintFull(os.Stdout)
	} else {
		metadata.Print(os.Stdout)
	}
	slog.Info("command done")
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolP("verbose", "v", false, "Show full version information")
}
