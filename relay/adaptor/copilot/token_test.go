package copilot

import (
	"os"
	"testing"
)

func TestGetToken(t *testing.T) {
	gt, _ := os.LookupEnv("GITHUB_TOKEN")
	token, err := GetToken(gt)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("token:", token)
}
