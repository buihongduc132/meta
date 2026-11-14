package connector

import "sync"

// reregisterEnv triggers a one-shot E2EE re-registration (icdc_register) of the
// EXISTING whatsmeow device on the next connect. It re-publishes the same
// identity key + signed prekey + device list; it never creates a new identity.
// Used to nudge peers holding a stale view of our device (cs-ks6.6).
const reregisterEnv = "MAUTRIX_META_E2EE_REREGISTER"

type reregisterGate struct {
	once sync.Once
}

// take reports true at most once per process, only for an existing device and
// only when the env flag is "1" — no retry loop on failure.
func (g *reregisterGate) take(isNew bool, flag string) bool {
	if isNew || flag != "1" {
		return false
	}
	fired := false
	g.once.Do(func() { fired = true })
	return fired
}

var e2eeReregister reregisterGate
