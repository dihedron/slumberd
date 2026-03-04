/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

// poweroffCmd represents the poweroff command
var poweroffCmd = &cobra.Command{
	Use:   "poweroff",
	Short: "Power off the system",
	Long: `Power off the system.

This command will immediately power off the system.`,
	Run: RunPoweroff,
}

func RunPoweroff(cmd *cobra.Command, args []string) {
	slog.Info("poweroff called")
	fmt.Println("system shutting down now...")
	//power.Shutdown()
	slog.Info("poweroff done")
}

func init() {
	rootCmd.AddCommand(poweroffCmd)
}
