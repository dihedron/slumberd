package command

import (
	"github.com/dihedron/devws/command/api"
	"github.com/dihedron/devws/command/power"
	"github.com/dihedron/devws/command/version"
)

// Commands is the set of root command groups.
type Commands struct {
	// API is the command that starts the API server.
	API api.API `command:"api" alias:"a" description:"Start the API server." `
	// Shutdown is the command that shuts down the machine.
	Shutdown power.Shutdown `command:"shutdown" alias:"s" description:"Shut down the machine."`
	// Hibernate is the command that hibernates the machine.
	Hibernate power.Hibernate `command:"hibernate" alias:"h" description:"Hibernate the machine."`
	// Version prints overlay version information and exits.
	Version version.Version `command:"version" alias:"v" description:"Show the command version and exit."`
}
