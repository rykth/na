//go:build !linux && !darwin

package unsupported

import (
	"fmt"
	"runtime"

	"github.com/rykth/na/internal/collector"
	"github.com/rykth/na/internal/config"
)

func init() {
	collector.Register("unsupported", newUnsupported)
}

func newUnsupported(_ *config.Config) (collector.Collector, error) {
	return nil, fmt.Errorf("na: unsupported OS: %s", runtime.GOOS)
}
