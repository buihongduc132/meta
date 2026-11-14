// bhd fork (cs-ks6.6 / cs-ks6.16): capture prekey messages that fail to decrypt,
// for offline analysis against the bridge's own key sets. Off unless
// WHATSMEOW_DUMP_FAILED_PKMSG_DIR is set. Writes CIPHERTEXT only (never plaintext),
// 0600, and logs public header fields as short prefixes.
package whatsmeow

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"go.mau.fi/libsignal/protocol"

	"go.mau.fi/whatsmeow/types"
)

const dumpDirEnv = "WHATSMEOW_DUMP_FAILED_PKMSG_DIR"

var unsafeFileChars = regexp.MustCompile(`[^A-Za-z0-9._-]`)

func shortHex(b []byte) string {
	if len(b) > 8 {
		b = b[:8]
	}
	return hex.EncodeToString(b)
}

// describePreKeyHeader returns the public header fields of a prekey message.
func describePreKeyHeader(content []byte) (string, error) {
	msg, err := protocol.NewPreKeySignalMessageFromBytes(content, pbSerializer.PreKeySignalMessage, pbSerializer.SignalMessage)
	if err != nil {
		return "", err
	}
	preKeyID := "none"
	if pk := msg.PreKeyID(); pk != nil && !pk.IsEmpty {
		preKeyID = fmt.Sprint(pk.Value)
	}
	return fmt.Sprintf("prekey_id=%s signed_prekey_id=%d sender_reg_id=%d base_key=%s sender_identity=%s version=%d",
		preKeyID, msg.SignedPreKeyID(), msg.RegistrationID(),
		shortHex(msg.BaseKey().Serialize()), shortHex(msg.IdentityKey().PublicKey().Serialize()), msg.MessageVersion()), nil
}

func (cli *Client) dumpFailedPreKeyMessage(info *types.MessageInfo, content any) {
	dir := os.Getenv(dumpDirEnv)
	raw, ok := content.([]byte)
	if dir == "" || !ok {
		return
	}
	header, err := describePreKeyHeader(raw)
	if err != nil {
		cli.Log.Warnf("[pkmsg-dump] %s: unparsable header: %v", info.ID, err)
		return
	}
	name := unsafeFileChars.ReplaceAllString(fmt.Sprintf("%s_%s.pkmsg", info.SourceString(), info.ID), "_")
	path := filepath.Join(dir, name)
	if err = os.MkdirAll(dir, 0o700); err == nil {
		err = os.WriteFile(path, raw, 0o600)
	}
	cli.Log.Warnf("[pkmsg-dump] %s from %s: %s file=%s write_err=%v", info.ID, info.SourceString(), header, path, err)
}
