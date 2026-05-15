//go:build linux

package linux

import "github.com/rickKoch/na/internal/collector"

func init() {
	collector.Register("linux", New)
}
