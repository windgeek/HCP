package manifest

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/windgeek/HCP/pkg/aha"
	"github.com/windgeek/HCP/pkg/zkp"
)

// Asset represents a single file in the release.
type Asset struct {
	Path      string `json:"path"`
	RawHash   string `json:"raw_hash"`
	LogicHash string `json:"logic_hash,omitempty"`
}

// Manifest represents the HCP Proof of Humanity.
type Manifest struct {
	Version         string                    `json:"version"`
	Author          string                    `json:"author"`       // Author's Address
	PublicKey       string                    `json:"public_key"`   // Hex encoded public key (added Phase 5)
	ContentHash     string                    `json:"content_hash"` // SHA256 of the content
	ParentHash      string                    `json:"parent_hash,omitempty"` // Provenance Chain (added Phase 6)
	Timestamp       int64                     `json:"timestamp"`
	EntropyDNA      string                    `json:"entropy_dna"`      // Random entropy for now
	Assets          []Asset                   `json:"assets,omitempty"` // Changed to []Asset in Phase 6
	ContributionMap map[string]aha.AHAMetrics `json:"contribution_map,omitempty"`
	CognitiveProofs map[string]zkp.Proof      `json:"cognitive_proofs,omitempty"` // Added Phase 4
	Signature       string                    `json:"signature"`                  // Hex encoded signature
}

// NewManifest creates a new Manifest for a given file.
func NewManifest(filePath string, authorAddr string, pubKey string) (*Manifest, error) {
	// 1. Calculate Content Hash
	hash, err := calculateFileHash(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %w", err)
	}

	// 2. Generate Random EntropyDNA
	entropy := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, entropy); err != nil {
		return nil, fmt.Errorf("failed to generate entropy: %w", err)
	}

	return &Manifest{
		Version:     "v1",
		Author:      authorAddr,
		PublicKey:   pubKey,
		ContentHash: hash,
		Timestamp:   time.Now().Unix(),
		EntropyDNA:  hex.EncodeToString(entropy),
	}, nil
}

// Sign signs the manifest using the provided private key.
// It signs the hash of the JSON representation (excluding the signature itself).
func (m *Manifest) Sign(key *btcec.PrivateKey) error {
	// 1. Serialize for signing (canonical JSON)
	// We need a stable representation. For simplicity, we create a struct without signature.
	type payload struct {
		Version         string                    `json:"version"`
		Author          string                    `json:"author"`
		PublicKey       string                    `json:"public_key"`
		ContentHash     string                    `json:"content_hash"`
		ParentHash      string                    `json:"parent_hash,omitempty"`
		Timestamp       int64                     `json:"timestamp"`
		EntropyDNA      string                    `json:"entropy_dna"`
		Assets          []Asset                   `json:"assets,omitempty"`
		ContributionMap map[string]aha.AHAMetrics `json:"contribution_map,omitempty"`
		CognitiveProofs map[string]zkp.Proof      `json:"cognitive_proofs,omitempty"`
	}
	p := payload{
		Version:         m.Version,
		Author:          m.Author,
		PublicKey:       m.PublicKey,
		ContentHash:     m.ContentHash,
		ParentHash:      m.ParentHash,
		Timestamp:       m.Timestamp,
		EntropyDNA:      m.EntropyDNA,
		Assets:          m.Assets,
		ContributionMap: m.ContributionMap,
		CognitiveProofs: m.CognitiveProofs,
	}

	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 2. Hash the data
	hash := sha256.Sum256(data)

	// 3. Sign
	signature := ecdsa.Sign(key, hash[:])

	// 4. Store signature
	m.Signature = hex.EncodeToString(signature.Serialize())
	return nil
}

// Save saves the signed manifest to a file.
func (m *Manifest) Save(path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func calculateFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
