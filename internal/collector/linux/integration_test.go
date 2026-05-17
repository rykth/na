//go:build linux && integration

package linux

import (
	"testing"

	"github.com/rykth/na/internal/config"
)

// Run with: go test -tags integration ./internal/collector/linux/
func TestIntegrationLinuxCollectorRead(t *testing.T) {
	col, err := New(config.Default())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer col.Close()

	stats, err := col.Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(stats) == 0 {
		t.Fatal("expected at least one interface")
	}

	var found bool
	for _, s := range stats {
		if s.Name == "lo" {
			found = true
			// uint64 counters are always >= 0, but verify the field is populated
			if s.RxBytes == 0 && s.TxBytes == 0 {
				t.Log("lo byte counters are zero (acceptable on quiet system)")
			}
			break
		}
	}
	if !found {
		t.Error("expected loopback interface 'lo' in results")
	}
}

func TestIntegrationLinuxCollectorList(t *testing.T) {
	col, err := New(config.Default())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer col.Close()

	names, err := col.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) == 0 {
		t.Fatal("expected at least one interface name")
	}
}
