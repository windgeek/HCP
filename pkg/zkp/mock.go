package zkp

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/windgeek/HCP/pkg/aha"
	"github.com/windgeek/HCP/pkg/cognitive"
)

// Proof represents a Zero-Knowledge Proof of Cognitive Work.
// In this MVP, it's a hash commitment to the metrics, signed by the timestamp.
type Proof struct {
	ProofID   string `json:"proof_id"`   // Unique ID of the proof (Hash)
	Timestamp int64  `json:"timestamp"`  // When the proof was generated
	PublicInput string `json:"public_input"` // Summary of what is being proven (e.g. "Complexity > 5")
}

// GenerateProof creates a mock ZKP for the given file's metrics.
// It proves:
// 1. The user spent time (proxied by AHA Commits)
// 2. The code has complexity (Cognitive Load)
// 3. The ratio is "Human"
func GenerateProof(metrics aha.AHAMetrics, complexity *cognitive.ComplexityStats) (*Proof, error) {
	// 1. Calculate Cognitive Ratio
	// Ratio = Complexity / (Commits + 1)
	// (Avoid division by zero)
	ratio := float64(complexity.Cyclomatic) / float64(metrics.Commits + 1)

	// 2. Construct the "Secret" Witness
	// In a real ZKP, this would be the private input to the circuit.
	witness := fmt.Sprintf("commits=%d;complexity=%d;volume=%f;ratio=%f;salt=%d",
		metrics.Commits, complexity.Cyclomatic, complexity.HalsteadVolume, ratio, time.Now().UnixNano())

	// 3. Generate the Proof (Hash Commitment)
	hasher := sha256.New()
	hasher.Write([]byte(witness))
	proofHash := hex.EncodeToString(hasher.Sum(nil))

	// 4. Public Input (What the verifier sees)
	// "I prove that complexity is N and commits are M, yielding a human ratio."
	publicInput := fmt.Sprintf("cyclomatic=%d;commits=%d", complexity.Cyclomatic, metrics.Commits)

	return &Proof{
		ProofID:     proofHash,
		Timestamp:   time.Now().Unix(),
		PublicInput: publicInput,
	}, nil
}

// VerifyProof checks if a proof is valid (Mock verification).
// Real ZKP would check the Groth16 proof here.
func VerifyProof(p *Proof) bool {
	// For MVP, just check if it has a valid hash structure
	if len(p.ProofID) != 64 {
		return false
	}
	return true
}
