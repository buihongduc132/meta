package whatsmeow

import (
	"os"
	"path/filepath"
	"testing"

	waLog "go.mau.fi/whatsmeow/util/log"

	"go.mau.fi/whatsmeow/types"
)

func TestDumpFailedPreKeyMessageOffByDefault(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(dumpDirEnv, "")
	cli := &Client{Log: waLog.Noop}
	cli.dumpFailedPreKeyMessage(&types.MessageInfo{ID: "m1"}, []byte{1, 2, 3})
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Fatal("must not write anything when the env flag is unset")
	}
}

func TestDumpFailedPreKeyMessageSkipsUnparsable(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "d")
	t.Setenv(dumpDirEnv, dir)
	cli := &Client{Log: waLog.Noop}
	cli.dumpFailedPreKeyMessage(&types.MessageInfo{ID: "m1"}, []byte{1, 2, 3})
	if _, err := os.Stat(dir); err == nil {
		t.Fatal("unparsable content must not be written")
	}
}
