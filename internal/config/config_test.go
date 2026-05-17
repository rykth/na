package config

import "testing"

func TestMatchInterfaceEmptyFilter(t *testing.T) {
	c := &Config{}
	for _, name := range []string{"eth0", "lo", "wlan0", "virbr0"} {
		if !c.MatchInterface(name) {
			t.Errorf("empty filter should match %q", name)
		}
	}
}

func TestMatchInterfaceIncludeGlob(t *testing.T) {
	c := &Config{Interfaces: "eth*"}
	if !c.MatchInterface("eth0") {
		t.Error("eth* should match eth0")
	}
	if !c.MatchInterface("eth1") {
		t.Error("eth* should match eth1")
	}
	if c.MatchInterface("wlan0") {
		t.Error("eth* should not match wlan0")
	}
	if c.MatchInterface("lo") {
		t.Error("eth* should not match lo")
	}
}

func TestMatchInterfaceExcludeGlob(t *testing.T) {
	c := &Config{Interfaces: "!lo"}
	if !c.MatchInterface("eth0") {
		t.Error("!lo should include eth0")
	}
	if !c.MatchInterface("wlan0") {
		t.Error("!lo should include wlan0")
	}
	if c.MatchInterface("lo") {
		t.Error("!lo should exclude lo")
	}
}

func TestMatchInterfaceIncludeExcludeMixed(t *testing.T) {
	c := &Config{Interfaces: "eth*,!eth2"}
	if !c.MatchInterface("eth0") {
		t.Error("eth*,!eth2 should match eth0")
	}
	if !c.MatchInterface("eth1") {
		t.Error("eth*,!eth2 should match eth1")
	}
	if c.MatchInterface("eth2") {
		t.Error("eth*,!eth2 should not match eth2")
	}
	if c.MatchInterface("wlan0") {
		t.Error("eth*,!eth2 should not match wlan0 (include pattern present)")
	}
}

func TestMatchInterfaceExactName(t *testing.T) {
	c := &Config{Interfaces: "eth0,lo"}
	if !c.MatchInterface("eth0") {
		t.Error("should match exact name eth0")
	}
	if !c.MatchInterface("lo") {
		t.Error("should match exact name lo")
	}
	if c.MatchInterface("eth1") {
		t.Error("should not match eth1 when not listed")
	}
}

func TestMatchInterfaceExcludeVirtual(t *testing.T) {
	c := &Config{Interfaces: "!virbr*,!docker*"}
	if !c.MatchInterface("eth0") {
		t.Error("should include eth0 when only excludes are present")
	}
	if c.MatchInterface("virbr0") {
		t.Error("should exclude virbr0")
	}
	if c.MatchInterface("docker0") {
		t.Error("should exclude docker0")
	}
}
