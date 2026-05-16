package model

import "testing"

func TestHistoryPushSamplesBasic(t *testing.T) {
	h := NewHistory(5)
	h.Push(1.0, 10.0)
	h.Push(2.0, 20.0)
	h.Push(3.0, 30.0)

	rx, tx := h.Samples(3)
	want := []float64{1.0, 2.0, 3.0}
	for i, v := range want {
		if rx[i] != v {
			t.Errorf("rx[%d]: expected %f, got %f", i, v, rx[i])
		}
		if tx[i] != v*10 {
			t.Errorf("tx[%d]: expected %f, got %f", i, v*10, tx[i])
		}
	}
}

func TestHistorySamplesZeroPadding(t *testing.T) {
	h := NewHistory(10)
	h.Push(5.0, 50.0)
	h.Push(6.0, 60.0)

	rx, tx := h.Samples(5)
	// first 3 slots should be zero-padded
	for i := range 3 {
		if rx[i] != 0.0 {
			t.Errorf("rx[%d]: expected 0 padding, got %f", i, rx[i])
		}
		if tx[i] != 0.0 {
			t.Errorf("tx[%d]: expected 0 padding, got %f", i, tx[i])
		}
	}
	if rx[3] != 5.0 {
		t.Errorf("rx[3]: expected 5.0, got %f", rx[3])
	}
	if rx[4] != 6.0 {
		t.Errorf("rx[4]: expected 6.0, got %f", rx[4])
	}
}

func TestHistoryWrapAround(t *testing.T) {
	h := NewHistory(3)
	h.Push(1.0, 0)
	h.Push(2.0, 0)
	h.Push(3.0, 0)
	// buffer is now full: [1, 2, 3], Head wraps to 0
	h.Push(4.0, 0)
	// oldest is now 2; buffer logically: [2, 3, 4]

	rx, _ := h.Samples(3)
	want := []float64{2.0, 3.0, 4.0}
	for i, v := range want {
		if rx[i] != v {
			t.Errorf("rx[%d]: expected %f, got %f", i, v, rx[i])
		}
	}
	if h.Count != 3 {
		t.Errorf("Count should stay at Size after wrap, got %d", h.Count)
	}
}

func TestHistoryFullBuffer(t *testing.T) {
	const size = 4
	h := NewHistory(size)
	for i := 1; i <= size; i++ {
		h.Push(float64(i), 0)
	}

	rx, _ := h.Samples(size)
	for i, want := range []float64{1, 2, 3, 4} {
		if rx[i] != want {
			t.Errorf("rx[%d]: expected %f, got %f", i, want, rx[i])
		}
	}
}

func TestHistorySamplesEmpty(t *testing.T) {
	h := NewHistory(10)
	rx, tx := h.Samples(5)
	for i, v := range rx {
		if v != 0 {
			t.Errorf("rx[%d]: expected 0, got %f", i, v)
		}
	}
	for i, v := range tx {
		if v != 0 {
			t.Errorf("tx[%d]: expected 0, got %f", i, v)
		}
	}
}

func TestHistorySamplesFewerThanRequested(t *testing.T) {
	h := NewHistory(10)
	h.Push(7.0, 70.0)

	rx, tx := h.Samples(3)
	if len(rx) != 3 || len(tx) != 3 {
		t.Fatalf("Samples(3) should return slices of length 3")
	}
	if rx[0] != 0 || rx[1] != 0 {
		t.Errorf("first two entries should be zero, got %v", rx[:2])
	}
	if rx[2] != 7.0 {
		t.Errorf("last entry should be 7.0, got %f", rx[2])
	}
}

func TestHistoryMultipleWraps(t *testing.T) {
	h := NewHistory(3)
	// push 7 samples through a buffer of size 3
	for i := 1; i <= 7; i++ {
		h.Push(float64(i), 0)
	}
	// last 3 pushed: 5, 6, 7
	rx, _ := h.Samples(3)
	want := []float64{5.0, 6.0, 7.0}
	for i, v := range want {
		if rx[i] != v {
			t.Errorf("rx[%d]: expected %f, got %f", i, v, rx[i])
		}
	}
}
