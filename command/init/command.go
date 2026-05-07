package init

import "github.com/dihedron/slumberd/command/base"

type Init struct {
	base.Command
}

func (cmd *Init) Execute(args []string) error {
	return nil
}
