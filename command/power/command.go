package power

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/dihedron/slumberd/command/base"
	"github.com/dihedron/slumberd/timex"
	"github.com/godbus/dbus/v5"
)

// Shutdown is the command that shuts down the machine.
type Shutdown struct {
	base.Command
	// Timeout is the duration to wait before shutting down.
	Timeout timex.Duration `short:"t" long:"timeout" description:"Timeout before shutting down." default:"0s"`
}

// Execute is the real implementation of the Shutdown command.
func (cmd *Shutdown) Execute(args []string) error {
	slog.Info("requesting system shutdown", "timeout", cmd.Timeout.String())
	if cmd.Timeout > 0 {
		time.Sleep(time.Duration(cmd.Timeout))
	}
	return callLogind("PowerOff")
}

// Hibernate is the command that hibernates the machine.
type Hibernate struct {
	base.Command
	// Timeout is the duration to wait before hibernating.
	Timeout timex.Duration `short:"t" long:"timeout" description:"Timeout before hibernating." default:"0s"`
}

// Execute is the real implementation of the Hibernate command.
func (cmd *Hibernate) Execute(args []string) error {
	slog.Info("requesting system hibernation", "timeout", cmd.Timeout.String())
	if cmd.Timeout > 0 {
		time.Sleep(time.Duration(cmd.Timeout))
	}
	return callLogind("Hibernate")
}

const (
	dbusDest      = "org.freedesktop.login1"
	dbusPath      = "/org/freedesktop/login1"
	dbusInterface = "org.freedesktop.login1.Manager"
)

// callLogind is a helper function that calls a method on the logind D-Bus service.
func callLogind(method string) error {
	conn, err := dbus.SystemBus()
	if err != nil {
		slog.Error("failed to connect to system bus", "error", err)
		return fmt.Errorf("failed to connect to system bus: %w", err)
	}
	slog.Debug("connected to system bus")
	defer conn.Close()

	obj := conn.Object(dbusDest, dbus.ObjectPath(dbusPath))

	// the boolean argument is for "interactive" (polkit dialog)
	if call := obj.Call(dbusInterface+"."+method, 0, true); call.Err != nil {
		slog.Error("dbus call failed", "method", method, "error", call.Err)
		return fmt.Errorf("dbus call to %s failed: %w", method, call.Err)
	}
	slog.Debug("dbus call successful")
	return nil
}
