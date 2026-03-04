/*
Copyright © 2026 Andrea Funtò dihedron.dev@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	packages  string
	debounce  time.Duration
	timeout   time.Duration
	frequency time.Duration
)

type Configuration struct {
	Packages  string        `json:"packages,omitempty" yaml:"packages,omitempty"`
	Debounce  time.Duration `json:"debounce,omitempty" yaml:"debounce,omitempty"`
	Timeout   time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Frequency time.Duration `json:"frequency,omitempty" yaml:"frequency,omitempty"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "slumberd",
	Short: "A daemon to manage system power state",
	Long:  `A daemon to manage system power state.`,
	Run:   RunRoot,
}

func RunRoot(cmd *cobra.Command, args []string) {
	slog.Info("running root command")
	fmt.Printf("packages: %s\n", packages)
	fmt.Printf("debounce: %s\n", debounce)
	fmt.Printf("timeout: %s\n", timeout)
	fmt.Printf("frequency: %s\n", frequency)

	slog.Info("command done")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/slumberd.yaml)")

	rootCmd.Flags().StringVar(&packages, "packages", "/home/developer/packages.yaml", "packages file")
	rootCmd.Flags().DurationVar(&debounce, "debounce", 500*time.Millisecond, "debounce time")
	rootCmd.Flags().DurationVar(&timeout, "timeout", 15*time.Minute, "timeout")
	rootCmd.Flags().DurationVar(&frequency, "frequency", time.Minute, "frequency")
	viper.BindPFlag("packages", rootCmd.Flags().Lookup("packages"))
	viper.BindPFlag("debounce", rootCmd.Flags().Lookup("debounce"))
	viper.BindPFlag("timeout", rootCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("frequency", rootCmd.Flags().Lookup("frequency"))
	viper.SetDefault("packages", "/home/developer/packages.yaml")
	viper.SetDefault("debounce", 500*time.Millisecond)
	viper.SetDefault("timeout", 15*time.Minute)
	viper.SetDefault("frequency", time.Minute)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".slumberd" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/")
		viper.SetConfigType("yaml")
		viper.SetConfigName("slumberd")
	}

	// Tells Viper to use this prefix when reading environment variables
	viper.SetEnvPrefix("SLUMBERD")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		slog.Info("using config file", "file", viper.ConfigFileUsed())
		fmt.Println("using config file", viper.ConfigFileUsed())
	}
}
