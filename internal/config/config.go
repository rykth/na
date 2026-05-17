package config

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

var Version = "dev" // overridden at build time via -X ldflags by GoReleaser

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

// Parse parses CLI flags into a Config and returns it.
// Prints version and exits on -v/--version; prints usage and exits on -h/--help.
func Parse() *Config {
	cfg := Default()
	var showVersion bool

	flag.StringVar(&cfg.Interfaces, "i", "", "comma-separated interface glob filter (e.g. \"eth*,!lo\")")
	flag.StringVar(&cfg.Interfaces, "interfaces", "", "comma-separated interface glob filter (e.g. \"eth*,!lo\")")
	flag.Float64Var(&cfg.Interval, "r", cfg.Interval, "poll interval in `seconds`")
	flag.Float64Var(&cfg.Interval, "interval", cfg.Interval, "poll interval in `seconds`")
	flag.BoolVar(&cfg.UseBits, "b", false, "display rates in bits per second")
	flag.BoolVar(&cfg.UseBits, "bits", false, "display rates in bits per second")
	flag.BoolVar(&cfg.UseSI, "s", false, "use SI units (kB/MB) instead of IEC (KiB/MiB)")
	flag.BoolVar(&cfg.UseSI, "si", false, "use SI units (kB/MB) instead of IEC (KiB/MiB)")
	flag.BoolVar(&cfg.ShowAll, "a", false, "show all interfaces including DOWN")
	flag.BoolVar(&cfg.ShowAll, "all", false, "show all interfaces including DOWN")
	flag.IntVar(&cfg.HistorySize, "history", cfg.HistorySize, "history ring buffer depth in samples")
	flag.BoolVar(&showVersion, "v", false, "print version and exit")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: na [OPTIONS]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -i, --interfaces string   interface glob filter (default: all)\n")
		fmt.Fprintf(os.Stderr, "  -r, --interval float      poll interval in seconds (default: 1.0)\n")
		fmt.Fprintf(os.Stderr, "  -b, --bits                display rates in bits per second\n")
		fmt.Fprintf(os.Stderr, "  -s, --si                  use SI units (kB/MB) instead of KiB/MiB\n")
		fmt.Fprintf(os.Stderr, "  -a, --all                 show all interfaces including DOWN\n")
		fmt.Fprintf(os.Stderr, "      --history int         history buffer depth (default: 120)\n")
		fmt.Fprintf(os.Stderr, "  -v, --version             print version and exit\n")
		fmt.Fprintf(os.Stderr, "  -h, --help                print this help\n")
	}

	flag.Parse()

	if showVersion {
		fmt.Printf("na %s\n", Version)
		os.Exit(0)
	}

	return cfg
}

// MatchInterface reports whether name passes the interface filter in cfg.Interfaces.
// Patterns are comma-separated globs; prefix ! to exclude.
// If any non-exclude pattern is present, unlisted interfaces are excluded by default.
func (c *Config) MatchInterface(name string) bool {
	if c.Interfaces == "" {
		return true
	}
	patterns := strings.Split(c.Interfaces, ",")

	hasInclude := false
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p != "" && !strings.HasPrefix(p, "!") {
			hasInclude = true
			break
		}
	}

	result := !hasInclude // default: included when only exclude patterns exist
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		exclude := strings.HasPrefix(p, "!")
		pat := p
		if exclude {
			pat = p[1:]
		}
		if ok, _ := path.Match(pat, name); ok {
			if exclude {
				return false
			}
			result = true
		}
	}
	return result
}

// LoadFile is a stub for future TOML config file support.
func (c *Config) LoadFile(_ string) error { return nil }
