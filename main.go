package main

import (
	"fmt"
	"os"

	"github.com/dihedron/slumberd/command"
	"github.com/jessevdk/go-flags"
)

func main() {
	defer cleanup()

	options := command.Commands{}
	if _, err := flags.NewParser(&options, flags.Default).Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		case *flags.Error:
			fmt.Fprintf(os.Stderr, "error: %s (%T)\n", err, err)
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
}

/*
func main() {
	defer cleanup()

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "v", "version", "-version", "--version":
			slog.Info("executing version")
			if len(os.Args) > 2 && (os.Args[2] == "--verbose" || os.Args[2] == "-v") {
				metadata.PrintFull(os.Stdout)
				os.Exit(0)
			} else {
				metadata.Print(os.Stdout)
				os.Exit(0)
			}
		case "p", "poweroff", "-poweroff", "--poweroff":
			slog.Info("executing poweroff")
			power.Shutdown()
			os.Exit(0)
		case "i", "init", "-init", "--init", "initialise", "-initialise", "--initialise", "g", "gen", "-gen", "--gen", "generate", "-generate", "--generate":
			slog.Info("executing init")
			cfg := configuration.Configuration{
				Packages:  pointer.To("/home/developer/packages.yaml"),
				Debounce:  pointer.To(timex.Duration(500 * time.Millisecond)),
				Timeout:   pointer.To(timex.Duration(15 * time.Minute)),
				Frequency: pointer.To(timex.Duration(time.Minute)),
			}
			if data, err := yaml.Marshal(cfg); err != nil {
				slog.Error("failed to marshal configuration", "error", err)
				fmt.Fprintf(os.Stderr, "failed to marshal configuration: %v\n", err)
				os.Exit(1)
			} else {
				fmt.Print(string(data))
			}
			os.Exit(0)
		}
	}

	var command Command

	var parser = flags.NewParser(&command, flags.Default)

	if args, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			//fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else {
		if err := command.Execute(args); err != nil {
			//fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}

}
*/
