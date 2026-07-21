package main

import (
	"os"
	"strings"
	"testing"
)

func TestMacBundleMetadataKeepsStableUserProfile(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("build/darwin/Info.plist")
	if err != nil {
		t.Fatalf("read macOS Info.plist: %v", err)
	}
	metadata := string(data)

	for _, expected := range []string{
		"<string>dev.tingly.tingly-ds</string>",
		"<string>tingly-ds</string>",
	} {
		if !strings.Contains(metadata, expected) {
			t.Errorf("Info.plist does not contain stable identity %q", expected)
		}
	}
	if strings.Contains(metadata, "<key>LSUIElement</key>") {
		t.Error("Info.plist enables LSUIElement, which hides a normal app from the Dock")
	}
}
