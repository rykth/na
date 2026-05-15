//go:build linux

package linux

import (
	"os"
	"strconv"
	"strings"
)

const sysfsNet = "/sys/class/net"

// ifaceMetadata reads per-interface metadata from /sys/class/net/<name>/.
// Missing or unreadable files are silently replaced with zero/empty values.
func ifaceMetadata(name string) (operState string, speed int, duplex string, mtu int, mac string, ifType int) {
	base := sysfsNet + "/" + name

	operState = sysfsString(base, "operstate")
	if operState == "" {
		operState = "unknown"
	}

	speed = sysfsInt(base, "speed", -1)
	duplex = sysfsString(base, "duplex")
	mtu = sysfsInt(base, "mtu", 0)
	mac = sysfsString(base, "address")
	ifType = sysfsInt(base, "type", 0)
	return
}

// listInterfaces returns all interface names found in /sys/class/net/.
func listInterfaces() ([]string, error) {
	entries, err := os.ReadDir(sysfsNet)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names, nil
}

func sysfsString(base, file string) string {
	data, err := os.ReadFile(base + "/" + file)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func sysfsInt(base, file string, fallback int) int {
	s := sysfsString(base, file)
	if s == "" {
		return fallback
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return v
}
