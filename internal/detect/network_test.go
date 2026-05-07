package detect

import (
	"testing"
	"time"
)

func TestHasActiveSSHConnections(t *testing.T) {
	for range 250 {
		time.Sleep(100 * time.Millisecond)
		active, err := HasActiveSSHConnections()
		if err != nil {
			t.Errorf("error checking for active SSH connections: %v", err)
		}
		if active {
			t.Log("active SSH connections found")
		} else {
			t.Log("no active SSH connections found")
		}
	}
}
