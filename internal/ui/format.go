package ui

import "fmt"

// formatRate formats a bytes/s rate according to the active display flags.
// useBits: multiply by 8 and display in bit/s units.
// useSI:   use 1000-based divisors (kB, MB, GB) instead of 1024-based (KiB, MiB, GiB).
func formatRate(bytesPerSec float64, useBits, useSI bool) string {
	if useBits {
		return formatScaled(bytesPerSec*8, 1000, "bit/s", "Kbit/s", "Mbit/s", "Gbit/s")
	}
	if useSI {
		return formatScaled(bytesPerSec, 1000, "B/s", "kB/s", "MB/s", "GB/s")
	}
	return formatScaled(bytesPerSec, 1024, "B/s", "KiB/s", "MiB/s", "GiB/s")
}

// formatBytes formats a cumulative byte total in human-readable form.
func formatBytes(total uint64, useSI bool) string {
	if useSI {
		return formatScaled(float64(total), 1000, "B", "kB", "MB", "GB")
	}
	return formatScaled(float64(total), 1024, "B", "KiB", "MiB", "GiB")
}

func formatScaled(v, div float64, unitB, unitK, unitM, unitG string) string {
	switch {
	case v >= div*div*div:
		return fmt.Sprintf("%.1f %s", v/(div*div*div), unitG)
	case v >= div*div:
		return fmt.Sprintf("%.1f %s", v/(div*div), unitM)
	case v >= div:
		return fmt.Sprintf("%.1f %s", v/div, unitK)
	default:
		return fmt.Sprintf("%.0f %s", v, unitB)
	}
}
