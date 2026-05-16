package model

import (
	"testing"
	"time"
)

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestCounterFirstSampleSkipsRate(t *testing.T) {
	var c Counter
	c.Update(1000, t0)
	if c.Rate != 0 {
		t.Errorf("first sample should have rate 0, got %f", c.Rate)
	}
	if c.Total != 1000 {
		t.Errorf("Total should be 1000, got %d", c.Total)
	}
}

func TestCounterRateCalculation(t *testing.T) {
	var c Counter
	c.Update(0, t0)
	c.Update(1000, t0.Add(time.Second))
	if c.Rate != 1000.0 {
		t.Errorf("expected rate 1000.0, got %f", c.Rate)
	}
}

func TestCounterRateWithFractionalInterval(t *testing.T) {
	var c Counter
	c.Update(0, t0)
	c.Update(500, t0.Add(500*time.Millisecond))
	if c.Rate != 1000.0 {
		t.Errorf("expected rate 1000.0, got %f", c.Rate)
	}
}

func TestCounterMonotonicDelta(t *testing.T) {
	// On 64-bit Linux kernels /proc/net/dev emits uint64 counters that never
	// realistically wrap. Verify that a large monotonic jump is measured
	// correctly.
	var c Counter
	const billion = uint64(1_000_000_000)
	c.Update(billion, t0)
	c.Update(billion+2000, t0.Add(time.Second))
	if c.Rate != 2000.0 {
		t.Errorf("expected rate 2000.0, got %f", c.Rate)
	}
}

func TestCounterReset(t *testing.T) {
	var c Counter
	c.Update(0, t0)
	c.Update(5000, t0.Add(time.Second))
	c.Reset()
	if c.Rate != 0 {
		t.Errorf("after Reset, Rate should be 0, got %f", c.Rate)
	}
	if c.Prev != c.Total {
		t.Errorf("after Reset, Prev should equal Total")
	}
	if !c.LastAt.IsZero() {
		t.Errorf("after Reset, LastAt should be zero")
	}
}

func TestCounterAccumulates(t *testing.T) {
	var c Counter
	c.Update(0, t0)
	c.Update(100, t0.Add(time.Second))
	c.Update(300, t0.Add(2*time.Second))
	// second interval: delta = 200 / 1s = 200
	if c.Rate != 200.0 {
		t.Errorf("expected rate 200.0, got %f", c.Rate)
	}
	if c.Total != 300 {
		t.Errorf("expected Total 300, got %d", c.Total)
	}
}
