package detect

import (
	"testing"
	"time"
)

func TestHasActiveSSHConnections(t *testing.T) {
	for range 100 {
		time.Sleep(250 * time.Millisecond)
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
