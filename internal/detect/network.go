package detect

import (
	"bufio"
	"log/slog"
	"os"
	"path"
	"strings"
)

var (
	tcpPath  = "/proc/net/tcp"
	tcp6Path = "/proc/net/tcp6"
)

// SetNetworkPaths allows overriding the paths to network proc files for testing.
func SetNetworkPaths(tcp, tcp6 string) {
	tcpPath = tcp
	tcp6Path = tcp6
}

// HasActiveSSHConnections checks if there are any established incoming SSH connections.
func HasActiveSSHConnections() (bool, error) {
	if active, err := checkTCP(tcpPath); err != nil {
		return false, err
	} else if active {
		return true, nil
	}
	return checkTCP(tcp6Path)
}

func checkTCP(filename string) (bool, error) {
	file, err := os.Open(path.Clean(filename))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Skip header
	if scanner.Scan() {
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}

			// Local address is field index 1
			// Remote address is field index 2
			// State is field index 3
			localAddr := fields[1]
			state := fields[3]

			// State "01" is ESTABLISHED
			if state != "01" {
				continue
			}

			// Check if local port is 22
			if strings.HasSuffix(localAddr, ":0016") {
				slog.Debug("active SSH connection found", "local", localAddr, "remote", fields[2])
				return true, nil
			}
		}
	}

	return false, scanner.Err()
}
