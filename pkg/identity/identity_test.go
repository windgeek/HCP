package identity

import (
	"path/filepath"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

func TestIdentityFlow(t *testing.T) {
	// 1. Generate Key
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	// 2. Derive Address
	addr, err := PubKeyToAddress(key.PubKey(), &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("PubKeyToAddress failed: %v", err)
	}
	if addr == "" {
		t.Fatal("Address is empty")
	}
	t.Logf("Generated Address: %s", addr)

	// 3. Save Key
	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "identity.key")
	passphrase := "securepassphrase"

	err = SaveKey(key, keyPath, passphrase)
	if err != nil {
		t.Fatalf("SaveKey failed: %v", err)
	}

	// 4. Load Key
	loadedKey, err := LoadKey(keyPath, passphrase)
	if err != nil {
		t.Fatalf("LoadKey failed: %v", err)
	}

	// 5. Verify Loaded Key matches Original
	if !key.PubKey().IsEqual(loadedKey.PubKey()) {
		t.Fatal("Loaded key pubkey does not match original")
	}

	// 6. Test Incorrect Passphrase
	_, err = LoadKey(keyPath, "wrongpassphrase")
	if err == nil {
		t.Fatal("LoadKey with wrong passphrase should fail")
	}
}
