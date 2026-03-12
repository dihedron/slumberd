package run

import (
	"log/slog"

	"github.com/dihedron/slumberd/command/base"
)

type Run struct {
	base.Command
}

func (cmd *Run) Execute(args []string) error {
	slog.Debug("running run command")
	return nil
}
