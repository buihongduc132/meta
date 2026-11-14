package connector

import "testing"

// cs-ks6.6 / @A D1: one-shot re-registration of an EXISTING e2ee device.
func TestReregisterGateFiresOnceForExistingDeviceWhenFlagged(t *testing.T) {
	var g reregisterGate
	if g.take(false, "") {
		t.Fatal("must not fire without the env flag")
	}
	if g.take(true, "1") {
		t.Fatal("must not fire for a brand-new device (it is registered anyway)")
	}
	if !g.take(false, "1") {
		t.Fatal("must fire for an existing device when flagged")
	}
	if g.take(false, "1") {
		t.Fatal("must fire at most once per process (no retry loop)")
	}
}
