# na — Network Analyzer

A keyboard-driven terminal UI for real-time network monitoring on Linux.

`na` reads live interface statistics directly from `/proc/net/dev` and `/sys/class/net/` — no root privileges required.

![CI](https://github.com/rykth/na/actions/workflows/ci.yml/badge.svg)

---

## Features

- Per-interface RX/TX bandwidth, packet rates, errors, and drops
- Scrolling ASCII sparkline graphs with sub-character precision (`▁▂▃▄▅▆▇█`)
- Scrollable, selectable interface list with live status indicators
- Configurable poll interval, unit display (IEC/SI/bits), and history depth
- Interface glob filter — include or exclude by name pattern
- Zero dependencies on external libraries beyond the TUI framework
- Runs as a regular user — no `sudo` needed

## Requirements

- Linux (reads `/proc/net/dev` and `/sys/class/net/`)
- Go 1.24+ (to build from source)

---

## Installation

### Pre-built binary (recommended)

Download the latest release for your architecture from the [Releases](https://github.com/rykth/na/releases) page:

```bash
# Linux amd64
curl -L https://github.com/rykth/na/releases/download/v0.1.0/na_0.1.0_linux_amd64.tar.gz | tar xz
sudo mv na /usr/local/bin/
```

```bash
# Linux arm64
curl -L https://github.com/rykth/na/releases/download/v0.1.0/na_0.1.0_linux_arm64.tar.gz | tar xz
sudo mv na /usr/local/bin/
```

Verify the download against the provided `checksums.txt`.

### From source

```bash
git clone https://github.com/rykth/na.git
cd na
make build          # produces bin/na
sudo mv bin/na /usr/local/bin/
```

Or install directly with Go:

```bash
go install github.com/rykth/na/cmd/na@latest
```

---

## Usage

```
na [OPTIONS]

Options:
  -i, --interfaces string   interface glob filter (default: all)
  -r, --interval float      poll interval in seconds (default: 1.0)
  -b, --bits                display rates in bits per second
  -s, --si                  use SI units (kB/MB) instead of KiB/MiB
  -a, --all                 show all interfaces including DOWN
      --history int         history buffer depth in samples (default: 120)
  -v, --version             print version and exit
  -h, --help                print this help
```

### Examples

```bash
# Watch all interfaces at 1 s interval
na

# Watch only Ethernet interfaces, exclude eth2
na -i 'eth*,!eth2'

# Show rates in Mbit/s, poll every 0.5 s
na -b -r 0.5

# Exclude virtual and loopback interfaces
na -i '!virbr*,!docker*,!lo'

# Keep 5 minutes of graph history at 1 s interval
na --history 300
```

---

## Keybindings

| Key | Action |
|-----|--------|
| `↑` / `k` | Select previous interface |
| `↓` / `j` | Select next interface |
| `g` | Cycle graph: off → RX → RX + TX |
| `d` | Toggle extended counters |
| `b` | Toggle bytes / bits display |
| `s` | Toggle SI / IEC units |
| `l` | Toggle interface list panel |
| `r` | Reset rate accumulators for selected interface |
| `+` / `-` | Increase / decrease poll interval (±0.5 s) |
| `h` / `?` | Toggle help overlay |
| `q` / `Ctrl-C` | Quit |

---

## Layout

```
┌─── na — network analyzer ────────────────────── 2025-05-17 12:34:56 ───┐
│ > eth0    ↓ 1.2 MiB/s ↑ 345 KiB/s [UP]  │ Interface: eth0              │
│   wlan0   ↓ 0 B/s     ↑ 0 B/s     [UP]  │ MAC: aa:bb:cc:dd:ee:ff       │
│   lo      ↓ 0 B/s     ↑ 0 B/s     [UP]  │ MTU: 1500  Speed: 1000 Mbps  │
│                                          │──────────────────────────────│
│                                          │ RX graph (toggle with g)     │
│                                          │ TX graph (toggle with g)     │
│                                          │──────────────────────────────│
│                                          │ RX/TX stats table            │
├──────────────────────────────────────────────────────────────────────────┤
│ interval: 1.0s                              q quit  ? help  ↑↓/jk select│
└──────────────────────────────────────────────────────────────────────────┘
```

---

## Building & Development

```bash
make build    # build to bin/na
make run      # go run ./cmd/na
make test     # go test ./...
make lint     # golangci-lint run
make clean    # remove bin/

# Integration tests (requires Linux, reads real /proc and /sys)
go test -tags integration ./internal/collector/linux/
```

---

## Data Sources

All data is read without elevated privileges:

| Source | Data |
|--------|------|
| `/proc/net/dev` | RX/TX bytes, packets, errors, drops, FIFO, collisions |
| `/sys/class/net/<iface>/operstate` | Link state (up/down/unknown) |
| `/sys/class/net/<iface>/speed` | Link speed in Mbps |
| `/sys/class/net/<iface>/duplex` | Duplex mode |
| `/sys/class/net/<iface>/mtu` | MTU |
| `/sys/class/net/<iface>/address` | MAC address |

---

## License

MIT
