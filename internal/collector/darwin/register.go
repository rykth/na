//go:build darwin

package darwin

import (
	"errors"

	"github.com/rykth/na/internal/collector"
	"github.com/rykth/na/internal/config"
)

func init() {
	collector.Register("darwin", newDarwin)
}

func newDarwin(_ *config.Config) (collector.Collector, error) {
	return nil, errors.New("na: macOS collector not yet implemented")
}
