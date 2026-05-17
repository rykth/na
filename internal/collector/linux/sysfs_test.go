//go:build linux

package linux

import (
	"os"
	"path/filepath"
	"testing"
)

func writeSysfsFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestSysfsString(t *testing.T) {
	dir := t.TempDir()
	writeSysfsFile(t, dir, "operstate", "up\n")

	if got := sysfsString(dir, "operstate"); got != "up" {
		t.Errorf("expected %q, got %q", "up", got)
	}
}

func TestSysfsStringMissing(t *testing.T) {
	dir := t.TempDir()
	if got := sysfsString(dir, "nonexistent"); got != "" {
		t.Errorf("expected empty string for missing file, got %q", got)
	}
}

func TestSysfsInt(t *testing.T) {
	dir := t.TempDir()
	writeSysfsFile(t, dir, "mtu", "1500\n")

	if got := sysfsInt(dir, "mtu", 0); got != 1500 {
		t.Errorf("expected 1500, got %d", got)
	}
}

func TestSysfsIntFallback(t *testing.T) {
	dir := t.TempDir()
	if got := sysfsInt(dir, "speed", -1); got != -1 {
		t.Errorf("expected fallback -1, got %d", got)
	}
}

func TestSysfsIntInvalid(t *testing.T) {
	dir := t.TempDir()
	writeSysfsFile(t, dir, "speed", "unknown\n")
	if got := sysfsInt(dir, "speed", -1); got != -1 {
		t.Errorf("expected fallback -1 for non-numeric value, got %d", got)
	}
}

func TestSysfsIntMultipleFiles(t *testing.T) {
	dir := t.TempDir()
	writeSysfsFile(t, dir, "type", "1\n")    // Ethernet
	writeSysfsFile(t, dir, "duplex", "full\n") // string, not int

	if got := sysfsInt(dir, "type", 0); got != 1 {
		t.Errorf("expected type 1, got %d", got)
	}
	if got := sysfsString(dir, "duplex"); got != "full" {
		t.Errorf("expected %q, got %q", "full", got)
	}
}
