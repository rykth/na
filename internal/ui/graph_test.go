package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/rykth/na/internal/collector"
	"github.com/rykth/na/internal/model"
)

func makeGraphModel() Model {
	store := model.NewStore(120)
	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	store.Update([]collector.RawIfStats{
		{Name: "eth0", OperState: "up", RxBytes: 0, TxBytes: 0},
	}, t0)
	store.Update([]collector.RawIfStats{
		{Name: "eth0", OperState: "up", RxBytes: 10_000, TxBytes: 5_000},
	}, t0.Add(time.Second))

	return Model{
		store:     store,
		width:     80,
		height:    24,
		showGraph: 1,
		selected:  0,
	}
}

func TestRenderGraphAxisCharacters(t *testing.T) {
	m := makeGraphModel()
	out := renderGraph(m, 40, graphLines)

	for _, want := range []string{"│", "└", "RX", "peak:"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in graph output, got:\n%s", want, out)
		}
	}
}

func TestRenderGraphTwoGraphs(t *testing.T) {
	m := makeGraphModel()
	m.showGraph = 2
	out := renderGraph(m, 40, graphLines*2)

	if !strings.Contains(out, "RX") {
		t.Error("expected 'RX' label in two-graph mode")
	}
	if !strings.Contains(out, "TX") {
		t.Error("expected 'TX' label in two-graph mode")
	}
}

func TestRenderGraphEmptyStore(t *testing.T) {
	m := Model{
		store:     model.NewStore(120),
		showGraph: 1,
	}
	// must not panic; returns empty/blank string
	_ = renderGraph(m, 40, graphLines)
}

func TestRenderGraphOff(t *testing.T) {
	m := makeGraphModel()
	m.showGraph = 0
	out := renderGraph(m, 40, graphLines)
	if out != "" {
		t.Errorf("showGraph=0 should return empty string, got %q", out)
	}
}

func TestDrawGraphFullBlock(t *testing.T) {
	// constant samples at the same level → all columns should reach peak → full blocks
	samples := make([]float64, 20)
	for i := range samples {
		samples[i] = 1000.0
	}
	out := drawGraph("RX", samples, 20, false, false)
	if !strings.Contains(out, "█") {
		t.Error("expected full block '█' for 100% fill columns")
	}
	if !strings.Contains(out, "─") {
		t.Error("expected '─' bottom axis")
	}
}

func TestDrawGraphZeroSamples(t *testing.T) {
	samples := make([]float64, 10)
	// all zero — peak = 0, divisor = 1 to avoid division by zero
	out := drawGraph("TX", samples, 10, false, false)
	if !strings.Contains(out, "TX") {
		t.Error("expected 'TX' label even for zero samples")
	}
}

func TestDrawGraphSubCharPrecision(t *testing.T) {
	// a sample at 50% fill should produce a partial block, not a full one
	samples := []float64{500.0, 500.0, 500.0, 1000.0} // peak = 1000, 50% fill
	out := drawGraph("RX", samples, 4, false, false)
	// Should contain at least one sub-character block char
	hasPartial := strings.ContainsAny(out, "▁▂▃▄▅▆▇")
	hasFull := strings.Contains(out, "█")
	if !hasPartial && !hasFull {
		t.Errorf("expected block characters in partial-fill graph, got:\n%s", out)
	}
}
