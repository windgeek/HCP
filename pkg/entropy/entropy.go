package entropy

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"time"
)

// Keystroke represents a single key press event.
type Keystroke struct {
	Key       rune      `json:"key"`
	Timestamp time.Time `json:"timestamp"`
}

// SessionStats contains the analysis of a typing session.
type SessionStats struct {
	TotalKeystrokes int     `json:"total_keystrokes"`
	Duration        float64 `json:"duration_seconds"` // Seconds
	MeanFlightTime  float64 `json:"mean_flight_time_ms"`
	Variance        float64 `json:"variance"`
	ShannonEntropy  float64 `json:"shannon_entropy"`
	HumanScore      float64 `json:"human_score"` // 0.0 to 1.0
}

// Analyze processes a sequence of keystrokes to extract entropy metrics.
func Analyze(keystrokes []Keystroke) (*SessionStats, error) {
	if len(keystrokes) < 2 {
		return nil, fmt.Errorf("not enough keystrokes to analyze")
	}

	start := keystrokes[0].Timestamp
	end := keystrokes[len(keystrokes)-1].Timestamp
	duration := end.Sub(start).Seconds()

	var flightTimes []float64
	var totalFlightTime float64

	// Calculate flight times (inter-key intervals)
	for i := 1; i < len(keystrokes); i++ {
		delta := keystrokes[i].Timestamp.Sub(keystrokes[i-1].Timestamp).Seconds() * 1000 // ms
		// Filter out unrealistic pauses (> 2 seconds is likely a break, not typing rhythm)
		if delta < 2000 {
			flightTimes = append(flightTimes, delta)
			totalFlightTime += delta
		}
	}

	count := float64(len(flightTimes))
	if count == 0 {
		return &SessionStats{TotalKeystrokes: len(keystrokes)}, nil
	}

	mean := totalFlightTime / count

	// Calculate Variance
	var sumSquaredDiff float64
	for _, t := range flightTimes {
		diff := t - mean
		sumSquaredDiff += diff * diff
	}
	variance := sumSquaredDiff / count

	// Calculate Shannon Entropy of the intervals (binned)
	// We use 50ms bins for the histogram
	bins := make(map[int]int)
	for _, t := range flightTimes {
		binIndex := int(t / 50)
		bins[binIndex]++
	}

	var shannon float64
	for _, binCount := range bins {
		p := float64(binCount) / count
		if p > 0 {
			shannon -= p * math.Log2(p)
		}
	}

	// Simple heuristic for "Human Score"
	// Humans tend to have some variance (not robotic 0) but not completely random.
	// A variance between 100 and 10000 ms^2 is typical for regular typing.
	// Shannon entropy usually > 1.5 bits for natural language typing.
	score := 0.0
	if variance > 50 && variance < 20000 {
		score += 0.5
	}
	if shannon > 1.0 {
		score += 0.5
	}

	return &SessionStats{
		TotalKeystrokes: len(keystrokes),
		Duration:        duration,
		MeanFlightTime:  mean,
		Variance:        variance,
		ShannonEntropy:  shannon,
		HumanScore:      score,
	}, nil
}

// GenerateSessionHash creates a unique hash of the session's biological markers.
// RFC-002: Hash(TimeDelta_1 || ... || TimeDelta_n)
func GenerateSessionHash(keystrokes []Keystroke) string {
	h := sha256.New()
	for i := 1; i < len(keystrokes); i++ {
		delta := keystrokes[i].Timestamp.Sub(keystrokes[i-1].Timestamp).Nanoseconds()
		h.Write([]byte(fmt.Sprintf("%d,", delta)))
	}
	return hex.EncodeToString(h.Sum(nil))
}
