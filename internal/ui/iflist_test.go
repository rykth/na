package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/rykth/na/internal/collector"
	"github.com/rykth/na/internal/model"
)

var t0ui = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func makeIfListModel() Model {
	store := model.NewStore(120)
	store.Update([]collector.RawIfStats{
		{Name: "eth0", OperState: "up", Speed: 1000, MTU: 1500},
		{Name: "wlan0", OperState: "down", Speed: -1, MTU: 1500},
		{Name: "lo", OperState: "up", Speed: -1, MTU: 65536},
	}, t0ui)
	return Model{
		store:    store,
		width:    80,
		height:   24,
		showList: true,
	}
}

func TestRenderIfListSortOrder(t *testing.T) {
	m := makeIfListModel()
	out := renderIfList(m, 20)

	eth0Idx := strings.Index(out, "eth0")
	wlan0Idx := strings.Index(out, "wlan0")
	loIdx := strings.Index(out, "lo")

	if eth0Idx < 0 || wlan0Idx < 0 || loIdx < 0 {
		t.Fatalf("not all interfaces found in output:\n%s", out)
	}
	if eth0Idx > wlan0Idx {
		t.Error("eth0 should appear before wlan0 (alphabetical)")
	}
	if loIdx < wlan0Idx {
		t.Error("lo should appear after wlan0 (lo pinned last)")
	}
}

func TestRenderIfListSelectionMarker(t *testing.T) {
	m := makeIfListModel()
	m.selected = 0
	out := renderIfList(m, 10)
	if !strings.Contains(out, ">") {
		t.Errorf("expected '>' selection marker, got:\n%s", out)
	}
}

func TestRenderIfListStateLabels(t *testing.T) {
	m := makeIfListModel()
	out := renderIfList(m, 10)
	if !strings.Contains(out, "[UP]") {
		t.Errorf("expected [UP] label, got:\n%s", out)
	}
	if !strings.Contains(out, "[DOWN]") {
		t.Errorf("expected [DOWN] label, got:\n%s", out)
	}
}

func TestRenderIfListScrolling(t *testing.T) {
	// with height=1 only the selected interface should be visible
	m := makeIfListModel()
	m.selected = 1 // wlan0
	out := renderIfList(m, 1)
	if !strings.Contains(out, "wlan0") {
		t.Errorf("selected interface should be visible in scrolled view, got:\n%s", out)
	}
	if strings.Contains(out, "eth0") {
		t.Errorf("non-selected interface should be scrolled out of view, got:\n%s", out)
	}
}

func TestRenderIfListRateDisplay(t *testing.T) {
	store := model.NewStore(120)
	// two updates so a rate is computed
	store.Update([]collector.RawIfStats{
		{Name: "eth0", OperState: "up", RxBytes: 0},
	}, t0ui)
	store.Update([]collector.RawIfStats{
		{Name: "eth0", OperState: "up", RxBytes: 1024},
	}, t0ui.Add(time.Second))

	m := Model{store: store, width: 80}
	out := renderIfList(m, 5)
	// rate should be non-zero; exact formatting depends on units
	if !strings.Contains(out, "↓") || !strings.Contains(out, "↑") {
		t.Errorf("expected rate arrows in output, got:\n%s", out)
	}
}
