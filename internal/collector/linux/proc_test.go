//go:build linux

package linux

import (
	"os"
	"strings"
	"testing"
)

func TestParseProcNetDev(t *testing.T) {
	f, err := os.Open("testdata/proc_net_dev.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	stats, err := parseProcNetDev(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stats) != 2 {
		t.Fatalf("expected 2 interfaces, got %d", len(stats))
	}

	lo := stats[0]
	if lo.Name != "lo" {
		t.Errorf("expected name %q, got %q", "lo", lo.Name)
	}
	if lo.RxBytes != 12345 {
		t.Errorf("lo RxBytes: expected 12345, got %d", lo.RxBytes)
	}
	if lo.TxPackets != 100 {
		t.Errorf("lo TxPackets: expected 100, got %d", lo.TxPackets)
	}

	eth := stats[1]
	if eth.Name != "eth0" {
		t.Errorf("expected name %q, got %q", "eth0", eth.Name)
	}
	if eth.RxBytes != 9876543 {
		t.Errorf("eth0 RxBytes: expected 9876543, got %d", eth.RxBytes)
	}
	if eth.RxErrors != 2 {
		t.Errorf("eth0 RxErrors: expected 2, got %d", eth.RxErrors)
	}
	if eth.TxColls != 9 {
		t.Errorf("eth0 TxColls: expected 9, got %d", eth.TxColls)
	}
	if eth.TxCompressed != 11 {
		t.Errorf("eth0 TxCompressed: expected 11, got %d", eth.TxCompressed)
	}
}

func TestParseProcNetDevTruncated(t *testing.T) {
	// only one header line — should return an error
	r := strings.NewReader("Inter-| Receive | Transmit\n")
	_, err := parseProcNetDev(r)
	if err == nil {
		t.Error("expected error for truncated header, got nil")
	}
}
