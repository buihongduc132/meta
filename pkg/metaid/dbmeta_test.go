package metaid

import (
	"testing"

	"go.mau.fi/mautrix-meta/pkg/messagix/types"
)

// Re-login (bridgev2 calls CopyFrom on the existing login) must keep the
// whatsmeow device: dropping WADeviceID makes connectE2EE register a brand-new
// E2EE device while Meta keeps publishing the old one, so peers encrypt to
// dead identities and inbound messages stop decrypting.
func TestCopyFromKeepsWADeviceIDOnMessengerRelogin(t *testing.T) {
	existing := &UserLoginMetadata{Platform: types.Facebook, WADeviceID: 155}
	fresh := &UserLoginMetadata{Platform: types.Facebook}
	existing.CopyFrom(fresh)
	if existing.WADeviceID != 155 {
		t.Fatalf("WADeviceID = %d after relogin, want 155 (device kept)", existing.WADeviceID)
	}
}

func TestCopyFromTakesNewWADeviceIDWhenProvided(t *testing.T) {
	existing := &UserLoginMetadata{Platform: types.Facebook, WADeviceID: 155}
	fresh := &UserLoginMetadata{Platform: types.Facebook, WADeviceID: 160}
	existing.CopyFrom(fresh)
	if existing.WADeviceID != 160 {
		t.Fatalf("WADeviceID = %d, want 160 (explicit new device wins)", existing.WADeviceID)
	}
}

func TestCopyFromDoesNotCarryDeviceAcrossNonMessengerPlatform(t *testing.T) {
	existing := &UserLoginMetadata{Platform: types.Facebook, WADeviceID: 155}
	fresh := &UserLoginMetadata{Platform: types.Instagram}
	existing.CopyFrom(fresh)
	if existing.WADeviceID != 0 {
		t.Fatalf("WADeviceID = %d, want 0 (no carry-over to Instagram)", existing.WADeviceID)
	}
}
