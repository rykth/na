package config

// Config holds application-wide settings populated from CLI flags.
type Config struct {
	Interfaces  string  // comma-separated glob filter, e.g. "eth*,!lo"
	Interval    float64 // poll interval in seconds
	UseBits     bool    // display rates in bits/s
	UseSI       bool    // use SI units (kB/MB) instead of KiB/MiB
	ShowAll     bool    // include DOWN interfaces
	HistorySize int     // history ring buffer depth
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		Interval:    1.0,
		HistorySize: 120,
	}
}
