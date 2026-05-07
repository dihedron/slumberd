package detect

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func IsAnyEditorActive2(procPath string) bool {

	pattern := regexp.MustCompile(`(.*)\.vscode-server\/cli\/servers\/.*\/server\/out\/bootstrap-fork.*--type=(?:fileWatcher|extensionHost)`)

	sshActive, err := HasActiveSSHConnections()
	if err != nil {
		slog.Error("failed to check SSH connections", "error", err)
		// If we can't check, we default to safety and assume it might be active
		// if the process is found. Or should we be strict?
		// Let's be strict as requested.
	}
	if !sshActive {
		slog.Debug("no active SSH connections, assuming editors are inactive/hung")
		fmt.Println(" > no active SSH connections, assuming editors are inactive/hung...")
		return false
	}

	files, err := os.ReadDir(procPath)
	if err != nil {
		slog.Error("failed to read /proc", "error", err)
		fmt.Println(" > failed to read /proc, assuming editors are inactive/hung...")
		return false
	}

	for _, f := range files {
		if !f.IsDir() || !isPID(f.Name()) {
			continue
		}

		filename := path.Clean(filepath.Join(procPath, f.Name(), "cmdline"))
		data, err := os.ReadFile(filename)
		if err != nil {
			continue
		}

		cmdline := strings.Replace(string(data), "\x00", " ", -1)
		if pattern.MatchString(cmdline) {
			slog.Debug("found active editor", "filename", filename, "cmdline", cmdline)
			return true
		}
	}

	slog.Debug("no active editors found")
	return false
}

var (
	editorRegex      = regexp.MustCompile(`vscode-server|code-server|cursor-server|windsurf-server|zed-remote-server|antigravity`)
	interpreterRegex = regexp.MustCompile(`^(node|python3?|sh|bash|perl|ruby)$`)
)

// IsAnyEditorActive checks if any of the target editor server components are running
// AND there is an active incoming SSH connection.
func IsAnyEditorActive(procPath string) []string {
	sshActive, err := HasActiveSSHConnections()
	if err != nil {
		slog.Error("failed to check SSH connections", "error", err)
		// If we can't check, we default to safety and assume it might be active
		// if the process is found. Or should we be strict?
		// Let's be strict as requested.
	}
	if !sshActive {
		slog.Debug("no active SSH connections, assuming editors are inactive/hung")
		fmt.Println(" > no active SSH connections, assuming editors are inactive/hung...")
		return nil
	}

	files, err := os.ReadDir(procPath)
	if err != nil {
		slog.Error("failed to read /proc", "error", err)
		fmt.Println(" > failed to read /proc, assuming editors are inactive/hung...")
		return nil
	}

	found := make(map[string]struct{})
	for _, f := range files {
		if !f.IsDir() || !isPID(f.Name()) {
			continue
		}

		cmdlinePath := filepath.Join(procPath, f.Name(), "cmdline")
		data, err := os.ReadFile(path.Clean(cmdlinePath))
		if err != nil {
			continue
		}

		parts := bytes.Split(data, []byte{0})
		if len(parts) == 0 {
			continue
		}

		// Check the executable itself
		argv0 := string(parts[0])
		if match := editorRegex.FindString(argv0); match != "" {
			found[match] = struct{}{}
			continue
		}

		// If it's a known interpreter, check for the script/program in arguments
		if interpreterRegex.MatchString(filepath.Base(argv0)) {
			for i := 1; i < len(parts); i++ {
				arg := string(parts[i])
				if len(arg) == 0 || strings.HasPrefix(arg, "-") {
					continue
				}
				// Found the first non-flag argument, check if it's our editor
				if match := editorRegex.FindString(arg); match != "" {
					found[match] = struct{}{}
				}
				break // Only consider the first non-flag argument as the "main thing"
			}
		}
	}

	if len(found) > 0 {
		var list []string
		for k := range found {
			list = append(list, k)
		}
		fmt.Printf(" > found active editors: %v\n", list)
		return list
	}

	fmt.Println(" > no active editors found...")
	return nil
}
