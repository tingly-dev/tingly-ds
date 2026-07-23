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

func TestWindowsMetadataIdentifiesUnprivilegedApplication(t *testing.T) {
	t.Parallel()

	manifest := readPackagingFile(t, "build/windows/wails.exe.manifest")
	for _, expected := range []string{
		`name="dev.tingly.tingly-ds"`,
		`level="asInvoker"`,
		`uiAccess="false"`,
		">permonitorv2,permonitor<",
		">UTF-8<",
	} {
		if !strings.Contains(manifest, expected) {
			t.Errorf("Windows manifest does not contain %q", expected)
		}
	}
	if strings.Contains(manifest, `level="requireAdministrator"`) {
		t.Error("Windows manifest requests administrator privileges")
	}

	versionInfo := readPackagingFile(t, "build/windows/info.json")
	for _, expected := range []string{`"ProductName": "Tingly DS"`, `"FileDescription": "Lightweight DeepSeek desktop shell"`} {
		if !strings.Contains(versionInfo, expected) {
			t.Errorf("Windows version metadata does not contain %q", expected)
		}
	}
}

func TestLinuxDesktopEntryKeepsStableIdentity(t *testing.T) {
	t.Parallel()

	desktop := readPackagingFile(t, "build/linux/tingly-ds.desktop")
	for _, expected := range []string{
		"Type=Application",
		"Name=Tingly DS",
		"Exec=tingly-ds",
		"Icon=tingly-ds",
		"Categories=Network;InstantMessaging;",
		"Terminal=false",
		"StartupNotify=true",
	} {
		if !strings.Contains(desktop, expected) {
			t.Errorf("Linux desktop entry does not contain %q", expected)
		}
	}
}

func readPackagingFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
