package model

// History is a fixed-size circular buffer of RX and TX rate samples.
type History struct {
	RX    []float64
	TX    []float64
	Size  int
	Head  int
	Count int
}

// NewHistory allocates a History with the given capacity.
func NewHistory(size int) *History {
	return &History{
		RX:   make([]float64, size),
		TX:   make([]float64, size),
		Size: size,
	}
}

// Push records a new RX/TX rate sample.
func (h *History) Push(rx, tx float64) {
	h.RX[h.Head] = rx
	h.TX[h.Head] = tx
	h.Head = (h.Head + 1) % h.Size
	if h.Count < h.Size {
		h.Count++
	}
}

// Samples returns the last n samples in chronological order (oldest→newest).
func (h *History) Samples(n int) (rx, tx []float64) {
	rx = make([]float64, n)
	tx = make([]float64, n)
	if h.Count == 0 {
		return
	}
	// number of valid samples we can actually return
	valid := min(h.Count, n)
	// start index: go back `valid` steps from Head
	start := (h.Head - valid + h.Size*2) % h.Size
	dst := n - valid // pad front with zeros
	for i := range valid {
		idx := (start + i) % h.Size
		rx[dst+i] = h.RX[idx]
		tx[dst+i] = h.TX[idx]
	}
	return
}
