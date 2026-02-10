package aha

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// AHAMetrics represents the Advanced Human Attribution scores for a file.
type AHAMetrics struct {
	Commits  int     `json:"commits"`   // Number of revisions
	AHAScore float64 `json:"aha_score"` // 0-100 score
}

// AnalyzeFile calculates the AHA metrics for a specific file.
func AnalyzeFile(filePath string, repoRoot string) (*AHAMetrics, error) {
	// Use relative path for git command
	relPath, err := filepath.Rel(repoRoot, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve relative path: %w", err)
	}

	// Git command to count commits affecting this file
	// --follow handles renames
	cmd := exec.Command("git", "log", "--follow", "--format=format:%H", "--", relPath)
	cmd.Dir = repoRoot
	output, err := cmd.Output()
	if err != nil {
		// If git fails (e.g. file not committed yet), return 0 commits
		// This is valid for new files.
		return &AHAMetrics{Commits: 0, AHAScore: 0}, nil
	}

	// Count lines (each line is a commit hash)
	commits := 0
	outStr := strings.TrimSpace(string(output))
	if outStr != "" {
		commits = len(strings.Split(outStr, "\n"))
	}

	// Calculate AHA Score
	// Simple Heuristic:
	// 1 commit = 10 (Initial commit)
	// 5 commits = 50 (Iterative refinement)
	// 10+ commits = 100 (Deep revisions)
	// New uncommitted files = 0
	
	score := float64(commits) * 10
	if score > 100 {
		score = 100
	}

	return &AHAMetrics{
		Commits:  commits,
		AHAScore: score,
	}, nil
}
