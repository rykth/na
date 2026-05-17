//go:build linux

package linux

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/rykth/na/internal/collector"
)

// parseProcNetDev parses the contents of /proc/net/dev from r.
func parseProcNetDev(r io.Reader) ([]collector.RawIfStats, error) {
	scanner := bufio.NewScanner(r)

	// skip the two header lines
	for range 2 {
		if !scanner.Scan() {
			return nil, fmt.Errorf("proc: /proc/net/dev: truncated header")
		}
	}

	var out []collector.RawIfStats
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		name, rest, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		name = strings.TrimSpace(name)
		fields := strings.Fields(rest)
		if len(fields) < 16 {
			continue
		}

		s := collector.RawIfStats{Name: name}
		s.RxBytes = parseU64(fields[0])
		s.RxPackets = parseU64(fields[1])
		s.RxErrors = parseU64(fields[2])
		s.RxDropped = parseU64(fields[3])
		s.RxFifo = parseU64(fields[4])
		s.RxFrame = parseU64(fields[5])
		s.RxCompressed = parseU64(fields[6])
		s.RxMulticast = parseU64(fields[7])
		s.TxBytes = parseU64(fields[8])
		s.TxPackets = parseU64(fields[9])
		s.TxErrors = parseU64(fields[10])
		s.TxDropped = parseU64(fields[11])
		s.TxFifo = parseU64(fields[12])
		s.TxColls = parseU64(fields[13])
		s.TxCarrier = parseU64(fields[14])
		s.TxCompressed = parseU64(fields[15])

		out = append(out, s)
	}
	return out, scanner.Err()
}

func readProcNetDev() ([]collector.RawIfStats, error) {
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseProcNetDev(f)
}

func parseU64(s string) uint64 {
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}
