package command

import (
	"github.com/dihedron/slumberd/command/version"
)

// Commands is the set of root command groups.
type Commands struct {
	// Version prints overlay version information and exits.
	Version version.Version `command:"version" alias:"v" description:"Show the command version and exit."`
}
