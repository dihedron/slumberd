package main

import (
	"log/slog"

	"github.com/dihedron/slumberd/version"
)

// Command is the main command that runs the application as
// a daemon (a systemd unit) in background.
type Command struct {

	// Version prints slumberd version information and exits.
	//lint:ignore SA5008 commands can have multiple aliases
	Version version.Version `command:"version" alias:"ver" alias:"v" description:"Show the command version and exit."`
}

// Execute runs the daemon command.
func (cmd *Command) Execute(args []string) error {
	slog.Info("starting daemon")
	return nil
}

/*
	slog.Info("starting daemon with configuration",
		"timeout", cmd.Configuration.Timeout,
		"frequency", cmd.Configuration.Frequency,
		"packages", *cmd.Configuration.Packages,
		"debounce", *cmd.Configuration.Debounce,
	)

	// set up signal handling for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	// set up ticker to run every frequency and check for active editors
	timeout := time.Duration(*cmd.Configuration.Timeout)
	frequency := time.Duration(*cmd.Configuration.Frequency)
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	// set up filesystem inotify watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("error setting up filesystem inotify watcher", "error", err)
		return err
	}
	if err := watcher.Add(filepath.Dir(*cmd.Configuration.Packages)); err != nil {
		slog.Error("error adding directory to filesystem inotify watcher", "path", filepath.Dir(*cmd.Configuration.Packages), "error", err)
		return err
	}
	defer watcher.Close()
	var timer *time.Timer
	var timerLock sync.Mutex

	lastActive := time.Now()

	for {
		select {
		case <-signals:
			slog.Info("received termination signal, shutting down")
			fmt.Println("received termination signal, shutting down...")
			return nil
		case event, ok := <-watcher.Events:
			if !ok {
				return fmt.Errorf("watcher closed")
			}
			slog.Info("event received", "event", event.Name, "operation", event.Op)
			if event.Name == *cmd.Configuration.Packages {
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					timerLock.Lock()
					// stop any existing timer (resetting the countdown)
					if timer != nil {
						timer.Stop()
					}
					// start a new timer
					timer = time.AfterFunc(time.Duration(*cmd.Configuration.Debounce), func() {
						slog.Info("file activity settled", "path", *cmd.Configuration.Packages)
						// TODO: read the file and install packages
					})
					timerLock.Unlock()
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return fmt.Errorf("watcher closed")
			}
			slog.Error("error from filesystem inotify watcher", "error", err)
		case <-ticker.C:
			editors := detect.IsAnyEditorActive2("/proc")
			if editors {
				slog.Info("editor sessions active")
				fmt.Println("editor sessions active...")
				lastActive = time.Now()
			} else {
				idleTime := time.Since(lastActive)
				slog.Info("no active editor sessions", "idle", idleTime.String())
				fmt.Printf("no active editor sessions... idle: %s\n", idleTime.String())
				if idleTime > timeout {
					slog.Warn("idle timeout reached, shutting down...")
					fmt.Println("shutting down...")
					//power.Shutdown()
					return nil
				}
			}
		}
	}
*/
