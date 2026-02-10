package identity

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

// GenerateKey generates a new secp256k1 private key.
func GenerateKey() (*btcec.PrivateKey, error) {
	return btcec.NewPrivateKey()
}

// PubKeyToAddress converts a public key to a Bitcoin P2WPKH (Native Segwit) address string.
// Supports MainNet and TestNet via chaincfg parameters.
func PubKeyToAddress(pubKey *btcec.PublicKey, net *chaincfg.Params) (string, error) {
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, net)
	if err != nil {
		return "", fmt.Errorf("failed to create address: %w", err)
	}
	return addr.EncodeAddress(), nil
}

// SaveKey encrypts and saves the private key to the specified path.
// The file format is hex(salt):hex(nonce):hex(ciphertext).
func SaveKey(key *btcec.PrivateKey, path string, passphrase string) error {
	// 1. Serialize the private key
	privBytes := key.Serialize()

	// 2. Generate a random salt
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// 3. Derive a key from the passphrase and salt
	encryptionKey := deriveKey(passphrase, salt)

	// 4. Create AES-GCM cipher
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	// 5. Generate a nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 6. Encrypt
	ciphertext := gcm.Seal(nil, nonce, privBytes, nil)

	// 7. Format data: salt:nonce:ciphertext (all hex encoded)
	data := fmt.Sprintf("%x:%x:%x", salt, nonce, ciphertext)

	// 8. Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 9. Write to file with restricted permissions
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	return nil
}

// LoadKey loads and decrypts the private key from the specified path.
func LoadKey(path string, passphrase string) (*btcec.PrivateKey, error) {
	// 1. Read file
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	// 2. Parse format
	parts := strings.Split(strings.TrimSpace(string(content)), ":")
	if len(parts) != 3 {
		return nil, errors.New("invalid key file format")
	}

	salt, err := hex.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid salt: %w", err)
	}

	nonce, err := hex.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid nonce: %w", err)
	}

	ciphertext, err := hex.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid ciphertext: %w", err)
	}

	// 3. Derive key
	encryptionKey := deriveKey(passphrase, salt)

	// 4. Create AES-GCM cipher
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// 5. Decrypt
	privBytes, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("decryption failed: invalid passphrase or corrupted file")
	}

	// 6. Parse private key
	privKey, _ := btcec.PrivKeyFromBytes(privBytes)
	return privKey, nil
}

// deriveKey derives a 32-byte key from a passphrase and salt using SHA-256.
// In a production environment, use a stronger KDF like Argon2 or Scrypt.
func deriveKey(passphrase string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(passphrase))
	hash.Write(salt)
	return hash.Sum(nil)
}
