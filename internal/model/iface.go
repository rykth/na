package model

import (
	"sort"
	"sync"
	"time"

	"github.com/rickKoch/na/internal/collector"
)

// Interface holds the live state of one network interface: metadata from
// sysfs, rate-aware counters derived from /proc/net/dev, and a rolling
// history of RX/TX rates used to render the sparkline graph.
type Interface struct {
	Name      string
	MAC       string
	Duplex    string
	OperState string // "up" | "down" | "unknown" | …
	Index     int
	MTU       int
	Speed     int // Mbps; -1 = unknown
	IfType    int // ARPHRD constant

	Stats   IfStats
	History *History
}

// Store is the central in-memory repository of all known interfaces.
// It is safe for concurrent access: the collector goroutine writes via
// Update() while the Bubbletea UI goroutine reads via Get() and Names().
type Store struct {
	mu          sync.RWMutex
	ifaces      map[string]*Interface
	historySize int
}

// NewStore creates an empty Store. historySize controls how many rate
// samples each interface's ring buffer can hold (passed through from Config).
func NewStore(historySize int) *Store {
	return &Store{
		ifaces:      make(map[string]*Interface),
		historySize: historySize,
	}
}

// Update ingests a fresh collector snapshot, advancing every counter and
// pushing new RX/TX rates into each interface's history ring buffer.
// New interfaces are created on first sight; existing ones are updated in
// place.
func (s *Store) Update(snapshots []collector.RawIfStats, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, snap := range snapshots {
		iface, exists := s.ifaces[snap.Name]
		if !exists {
			iface = &Interface{
				Name:    snap.Name,
				History: NewHistory(s.historySize),
			}
			s.ifaces[snap.Name] = iface
		}

		// refresh sysfs metadata every poll
		iface.OperState = snap.OperState
		iface.Speed = snap.Speed
		iface.Duplex = snap.Duplex
		iface.MTU = snap.MTU
		iface.MAC = snap.MAC
		iface.IfType = snap.IfType

		// advance all 16 rate counters
		iface.Stats.update(rawCounters{
			rxBytes: snap.RxBytes, rxPackets: snap.RxPackets,
			rxErrors: snap.RxErrors, rxDropped: snap.RxDropped,
			rxFifo: snap.RxFifo, rxFrame: snap.RxFrame,
			rxCompressed: snap.RxCompressed, rxMulticast: snap.RxMulticast,
			txBytes: snap.TxBytes, txPackets: snap.TxPackets,
			txErrors: snap.TxErrors, txDropped: snap.TxDropped,
			txFifo: snap.TxFifo, txColls: snap.TxColls,
			txCarrier: snap.TxCarrier, txCompressed: snap.TxCompressed,
		}, now)

		// record the new rates in the history ring buffer
		iface.History.Push(iface.Stats.RxBytes.Rate, iface.Stats.TxBytes.Rate)
	}
}

// Get returns the Interface for name, or nil if it is not known.
func (s *Store) Get(name string) *Interface {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ifaces[name]
}

// Names returns all known interface names in sorted order, with "lo" last.
func (s *Store) Names() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	names := make([]string, 0, len(s.ifaces))
	for name := range s.ifaces {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		if names[i] == "lo" {
			return false
		}
		if names[j] == "lo" {
			return true
		}
		return names[i] < names[j]
	})
	return names
}

// ResetStats zeroes the rate accumulators for the named interface so that
// rates restart from the next poll. Called when the user presses 'r'.
func (s *Store) ResetStats(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	iface, ok := s.ifaces[name]
	if !ok {
		return
	}
	st := &iface.Stats
	st.RxBytes.Reset()
	st.RxPackets.Reset()
	st.RxErrors.Reset()
	st.RxDropped.Reset()
	st.RxFifo.Reset()
	st.RxFrame.Reset()
	st.RxCompressed.Reset()
	st.RxMulticast.Reset()
	st.TxBytes.Reset()
	st.TxPackets.Reset()
	st.TxErrors.Reset()
	st.TxDropped.Reset()
	st.TxFifo.Reset()
	st.TxColls.Reset()
	st.TxCarrier.Reset()
	st.TxCompressed.Reset()
}
