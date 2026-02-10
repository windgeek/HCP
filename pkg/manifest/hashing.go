package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/windgeek/HCP/pkg/aha"
	"github.com/windgeek/HCP/pkg/cognitive"
	"github.com/windgeek/HCP/pkg/hash"
	"github.com/windgeek/HCP/pkg/zkp"
)

// CalculateDirHash scans a directory, ignores files, calculates global hash,
// and generates AHA/Cognitive metrics.
func CalculateDirHash(root string, ignorePatterns []string) (
	string, 
	[]Asset, // Changed return type
	map[string]aha.AHAMetrics, 
	map[string]zkp.Proof, 
	error,
) {
	var files []string
	
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		relPath, _ := filepath.Rel(root, path)
		if relPath == "." {
			return nil
		}

		if ShouldIgnore(relPath, ignorePatterns) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		
		// Additional hidden file check
		if strings.HasPrefix(filepath.Base(path), ".") && relPath != "." {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return "", nil, nil, nil, err
	}

	sort.Strings(files)

	globalHasher := sha256.New()
	var assets []Asset // Changed type
	contribMap := make(map[string]aha.AHAMetrics)
	zkpMap := make(map[string]zkp.Proof)
	
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return "", nil, nil, nil, err
		}
		
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			f.Close()
			return "", nil, nil, nil, err
		}
		f.Close()
		
		fileHash := hex.EncodeToString(h.Sum(nil))
		relPath, _ := filepath.Rel(root, file)
		cleanPath := filepath.ToSlash(relPath)

		// Calculate Logic Hash for .go files
		var logicHash string
		if strings.HasSuffix(file, ".go") {
			if lh, err := hash.ComputeLogicHash(file); err == nil {
				logicHash = lh
			}
		}

		asset := Asset{
			Path:      cleanPath,
			RawHash:   fileHash,
			LogicHash: logicHash,
		}
		assets = append(assets, asset)
		
		globalHasher.Write([]byte(cleanPath))
		globalHasher.Write([]byte(fileHash))
		// Note: We intentionally hash only RawHash into GlobalHash to maintain strict integrity chain.
		// LogicHash is for "Fuzzy Verification".

		// Analyze AHA
		metrics, err := aha.AnalyzeFile(file, root)
		if err != nil {
			// Warn mostly expected for non-git or new files
			metrics = &aha.AHAMetrics{Commits: 0, AHAScore: 0}
		}
		contribMap[cleanPath] = *metrics

		// Phase 4: Cognitive & ZKP
		complexity, err := cognitive.AnalyzeComplexity(file)
		if err == nil {
			proof, err := zkp.GenerateProof(*metrics, complexity)
			if err == nil {
				zkpMap[cleanPath] = *proof
			}
		}
	}

	return hex.EncodeToString(globalHasher.Sum(nil)), assets, contribMap, zkpMap, nil
}

// ShouldIgnore checks if a file path matches any ignore pattern.
func ShouldIgnore(path string, patterns []string) bool {
	for _, p := range patterns {
		matched, _ := filepath.Match(p, path)
		if matched {
			return true
		}
		// Also check partial path (directory ignore)
		if strings.Contains(path, "/"+p+"/") || strings.HasPrefix(path, p+"/") {
			return true
		}
		if p == filepath.Base(path) {
			return true
		}
	}
	return false
}
