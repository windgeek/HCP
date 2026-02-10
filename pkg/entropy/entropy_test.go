package entropy

import (
	"testing"
	"time"
)

func TestAnalyze(t *testing.T) {
	// Simulate typing "test" with varying intervals
	start := time.Now()
	keystrokes := []Keystroke{
		{Key: 't', Timestamp: start},
		{Key: 'e', Timestamp: start.Add(100 * time.Millisecond)},
		{Key: 's', Timestamp: start.Add(250 * time.Millisecond)}, // 150ms delta
		{Key: 't', Timestamp: start.Add(350 * time.Millisecond)}, // 100ms delta
	}

	stats, err := Analyze(keystrokes)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	if stats.TotalKeystrokes != 4 {
		t.Errorf("Expected 4 keystrokes, got %d", stats.TotalKeystrokes)
	}

	// Flight times: 100ms, 150ms, 100ms
	// Mean: (100+150+100)/3 = 116.66ms
	expectedMean := 116.66
	if stats.MeanFlightTime < expectedMean-1 || stats.MeanFlightTime > expectedMean+1 {
		t.Errorf("Mean flight time mismatch: got %f, want ~%f", stats.MeanFlightTime, expectedMean)
	}

	// Variance should be non-zero
	if stats.Variance == 0 {
		t.Error("Variance is 0, expected > 0")
	}

	t.Logf("Stats: %+v", stats)
}
