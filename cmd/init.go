/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a default configuration file",
	Long: `Generate a default configuration file.

This command will generate a default configuration file for the 
slumberd daemon; the output value can be redirected to a file on 
disk. Usually, the configuration file will be stored under
/etc/slumberd.yaml, in the current directory or in the user's home
directory.
`,
	Run: RunInit,
}

func RunInit(cmd *cobra.Command, args []string) {
	slog.Info("running init command")
	cfg := Configuration{
		Packages:  "/home/developer/packages.yaml",
		Debounce:  500 * time.Millisecond,
		Timeout:   15 * time.Minute,
		Frequency: time.Minute,
	}
	if data, err := yaml.Marshal(cfg); err != nil {
		slog.Error("failed to marshal configuration", "error", err)
		fmt.Fprintf(os.Stderr, "failed to marshal configuration: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Print(string(data))
	}
	slog.Info("command done")
}

func init() {
	rootCmd.AddCommand(initCmd)
}
