//go:build linux

package linux

import (
	"github.com/rickKoch/na/internal/collector"
	"github.com/rickKoch/na/internal/config"
)

// LinuxCollector reads network statistics from /proc and /sys on Linux.
type LinuxCollector struct {
	cfg *config.Config
}

// New returns a LinuxCollector.
func New(cfg *config.Config) (collector.Collector, error) {
	return &LinuxCollector{cfg: cfg}, nil
}

func (c *LinuxCollector) List() ([]string, error) {
	return listInterfaces()
}

// Read returns a snapshot of all interfaces with counters and sysfs metadata.
func (c *LinuxCollector) Read() ([]collector.RawIfStats, error) {
	stats, err := readProcNetDev()
	if err != nil {
		return nil, err
	}
	for i := range stats {
		name := stats[i].Name
		stats[i].OperState, stats[i].Speed, stats[i].Duplex,
			stats[i].MTU, stats[i].MAC, stats[i].IfType = ifaceMetadata(name)
	}
	return stats, nil
}

func (c *LinuxCollector) Close() error { return nil }
