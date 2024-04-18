package copilot

import (
	"os"
	"testing"
)

func TestGetMachineID(t *testing.T) {
	gt, _ := os.LookupEnv("GITHUB_TOKEN")
	machineID := GetMachineID(gt)
	t.Log(machineID)
}

func TestGetSessionID(t *testing.T) {
	gt, _ := os.LookupEnv("GITHUB_TOKEN")
	sessionID := GetSessionID(gt)
	t.Log(sessionID)
}
