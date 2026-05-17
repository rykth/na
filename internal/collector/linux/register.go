//go:build linux

package linux

import "github.com/rykth/na/internal/collector"

func init() {
	collector.Register("linux", New)
}
