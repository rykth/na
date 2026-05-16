package model

import (
	"time"
)

// Counter tracks a single cumulative kernel counter and derives a per-second
// rate.
type Counter struct {
	Total  uint64    // cumulative value from the kernel (monotonically inc.)
	Prev   uint64    // value at the previous poll
	Rate   float64   // (Total - Prev) / elapsed seconds
	LastAt time.Time // timestamp of the last Update call
}

// Update advances the counter with a new kernel snapshot taken at now.
// The first call initialises Prev and LastAt without computing a rate.
// On 64-bit Linux kernels /proc/net/dev emits uint64 counters; unsigned
// subtraction correctly handles any monotonic delta without overflow concerns.
func (c *Counter) Update(newTotal uint64, now time.Time) {
	if !c.LastAt.IsZero() {
		dt := now.Sub(c.LastAt).Seconds()
		if dt > 0 {
			delta := newTotal - c.Prev // unsigned wrap-around safe
			c.Rate = float64(delta) / dt
		}
	}
	c.Total = newTotal
	c.Prev = newTotal
	c.LastAt = now
}

// Reset zeroes the rate and restarts delta tracking from the current total.
// Called when the user presses 'r' to reset accumulators for an interface.
func (c *Counter) Reset() {
	c.Rate = 0
	c.Prev = c.Total
	c.LastAt = time.Time{}
}

// IfStats holds one Counter per field exposed by /proc/net/dev.
type IfStats struct {
	RxBytes      Counter
	RxPackets    Counter
	RxErrors     Counter
	RxDropped    Counter
	RxFifo       Counter
	RxFrame      Counter
	RxCompressed Counter
	RxMulticast  Counter

	TxBytes      Counter
	TxPackets    Counter
	TxErrors     Counter
	TxDropped    Counter
	TxFifo       Counter
	TxColls      Counter
	TxCarrier    Counter
	TxCompressed Counter
}

// update advances every counter in s from a raw snapshot taken at now.
func (s *IfStats) update(raw rawCounters, now time.Time) {
	s.RxBytes.Update(raw.rxBytes, now)
	s.RxPackets.Update(raw.rxPackets, now)
	s.RxErrors.Update(raw.rxErrors, now)
	s.RxDropped.Update(raw.rxDropped, now)
	s.RxFifo.Update(raw.rxFifo, now)
	s.RxFrame.Update(raw.rxFrame, now)
	s.RxCompressed.Update(raw.rxCompressed, now)
	s.RxMulticast.Update(raw.rxMulticast, now)

	s.TxBytes.Update(raw.txBytes, now)
	s.TxPackets.Update(raw.txPackets, now)
	s.TxErrors.Update(raw.txErrors, now)
	s.TxDropped.Update(raw.txDropped, now)
	s.TxFifo.Update(raw.txFifo, now)
	s.TxColls.Update(raw.txColls, now)
	s.TxCarrier.Update(raw.txCarrier, now)
	s.TxCompressed.Update(raw.txCompressed, now)
}

// rawCounters is an internal mirror of collector.RawIfStats counter fields,
// kept here to avoid importing the collector package from model.
type rawCounters struct {
	rxBytes, rxPackets, rxErrors, rxDropped    uint64
	rxFifo, rxFrame, rxCompressed, rxMulticast uint64
	txBytes, txPackets, txErrors, txDropped    uint64
	txFifo, txColls, txCarrier, txCompressed   uint64
}
