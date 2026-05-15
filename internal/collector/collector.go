package collector

import (
	"fmt"

	"github.com/rickKoch/na/internal/config"
)

// RawIfStats is a single kernel snapshot for one network interface.
// Counter fields come from /proc/net/dev; metadata fields from /sys/class/net.
type RawIfStats struct {
	Name string

	// /proc/net/dev RX fields (in order)
	RxBytes      uint64
	RxPackets    uint64
	RxErrors     uint64
	RxDropped    uint64
	RxFifo       uint64
	RxFrame      uint64
	RxCompressed uint64
	RxMulticast  uint64

	// /proc/net/dev TX fields (in order)
	TxBytes      uint64
	TxPackets    uint64
	TxErrors     uint64
	TxDropped    uint64
	TxFifo       uint64
	TxColls      uint64
	TxCarrier    uint64
	TxCompressed uint64

	// /sys/class/net/<iface>/ metadata
	OperState string // "up" | "down" | "unknown"
	Speed     int    // Mbps; -1 if unavailable
	Duplex    string // "full" | "half" | ""
	MTU       int
	MAC       string
	IfType    int // ARPHRD type: 1=Ethernet, 772=loopback
}

// Collector reads network interface statistics from the OS.
type Collector interface {
	List() ([]string, error)
	Read() ([]RawIfStats, error)
	Close() error
}

type factoryFn func(*config.Config) (Collector, error)

var registry = map[string]factoryFn{}

// Register associates a platform name with a Collector factory.
// Called from platform-specific init() functions.
func Register(name string, f factoryFn) {
	registry[name] = f
}

// New returns the Collector registered for the current platform.
func New(cfg *config.Config) (Collector, error) {
	for _, f := range registry {
		return f(cfg)
	}
	return nil, fmt.Errorf("collector: no implementation registered for this platform")
}
