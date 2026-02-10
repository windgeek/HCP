package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/windgeek/HCP/pkg/aha"
	"github.com/windgeek/HCP/pkg/zkp"
)

// Verify verifies the signature of the manifest against the provided public key.
func (m *Manifest) Verify(pubKey *btcec.PublicKey) error {
	// 1. Reconstruct payload (must match Sign method's payload)
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

	// 2. Hash
	hash := sha256.Sum256(data)

	// 3. Decode Signature
	sigBytes, err := hex.DecodeString(m.Signature)
	if err != nil {
		return fmt.Errorf("invalid signature format: %w", err)
	}
	signature, err := ecdsa.ParseSignature(sigBytes)
	if err != nil {
		return fmt.Errorf("failed to parse signature: %w", err)
	}

	// 4. Verify
	if !signature.Verify(hash[:], pubKey) {
		return fmt.Errorf("signature verification failed")
	}

	return nil
}
