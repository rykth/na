package model

import (
	"testing"
	"time"

	"github.com/rykth/na/internal/collector"
)

func makeSnap(name string, rxBytes, txBytes uint64) collector.RawIfStats {
	return collector.RawIfStats{
		Name:      name,
		RxBytes:   rxBytes,
		TxBytes:   txBytes,
		OperState: "up",
		Speed:     1000,
		MTU:       1500,
	}
}

func TestStoreUpdateCreatesInterfaces(t *testing.T) {
	store := NewStore(120)
	snaps := []collector.RawIfStats{
		makeSnap("eth0", 0, 0),
		makeSnap("lo", 0, 0),
	}
	store.Update(snaps, t0)

	names := store.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 interfaces, got %d", len(names))
	}
	// lo should be last
	if names[len(names)-1] != "lo" {
		t.Errorf("expected lo last, got %v", names)
	}
}

func TestStoreNamesLoLast(t *testing.T) {
	store := NewStore(120)
	snaps := []collector.RawIfStats{
		makeSnap("lo", 0, 0),
		makeSnap("eth0", 0, 0),
		makeSnap("wlan0", 0, 0),
	}
	store.Update(snaps, t0)

	names := store.Names()
	if names[len(names)-1] != "lo" {
		t.Errorf("lo should be last, got %v", names)
	}
	if names[0] != "eth0" {
		t.Errorf("eth0 should be first, got %v", names)
	}
}

func TestStoreUpdateRateCalculation(t *testing.T) {
	store := NewStore(120)
	store.Update([]collector.RawIfStats{makeSnap("eth0", 0, 0)}, t0)
	store.Update([]collector.RawIfStats{makeSnap("eth0", 1000, 500)}, t0.Add(time.Second))

	iface := store.Get("eth0")
	if iface == nil {
		t.Fatal("eth0 not found in store")
	}
	if iface.Stats.RxBytes.Rate != 1000.0 {
		t.Errorf("RxBytes.Rate: expected 1000.0, got %f", iface.Stats.RxBytes.Rate)
	}
	if iface.Stats.TxBytes.Rate != 500.0 {
		t.Errorf("TxBytes.Rate: expected 500.0, got %f", iface.Stats.TxBytes.Rate)
	}
}

func TestStoreGetUnknown(t *testing.T) {
	store := NewStore(120)
	if store.Get("notexist") != nil {
		t.Error("expected nil for unknown interface")
	}
}

func TestStoreMetadataRefreshed(t *testing.T) {
	store := NewStore(120)
	store.Update([]collector.RawIfStats{makeSnap("eth0", 0, 0)}, t0)

	snap2 := makeSnap("eth0", 0, 0)
	snap2.OperState = "down"
	snap2.Speed = -1
	store.Update([]collector.RawIfStats{snap2}, t0.Add(time.Second))

	iface := store.Get("eth0")
	if iface.OperState != "down" {
		t.Errorf("expected OperState 'down', got %q", iface.OperState)
	}
	if iface.Speed != -1 {
		t.Errorf("expected Speed -1, got %d", iface.Speed)
	}
}

func TestStoreResetStats(t *testing.T) {
	store := NewStore(120)
	store.Update([]collector.RawIfStats{makeSnap("eth0", 0, 0)}, t0)
	store.Update([]collector.RawIfStats{makeSnap("eth0", 9000, 3000)}, t0.Add(time.Second))

	store.ResetStats("eth0")
	iface := store.Get("eth0")
	if iface.Stats.RxBytes.Rate != 0 {
		t.Errorf("expected rate 0 after reset, got %f", iface.Stats.RxBytes.Rate)
	}
}

func TestStoreHistoryPushedOnUpdate(t *testing.T) {
	store := NewStore(10)
	store.Update([]collector.RawIfStats{makeSnap("eth0", 0, 0)}, t0)
	store.Update([]collector.RawIfStats{makeSnap("eth0", 2000, 1000)}, t0.Add(time.Second))

	iface := store.Get("eth0")
	if iface.History.Count < 1 {
		t.Errorf("expected at least 1 history sample, got %d", iface.History.Count)
	}
}
