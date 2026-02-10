package cognitive

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ComplexityStats holds the structural metrics of a file.
type ComplexityStats struct {
	NodeCount       int     `json:"node_count"`
	Cyclomatic      int     `json:"cyclomatic_complexity"`
	HalsteadVolume  float64 `json:"halstead_volume"`
	Functions       int     `json:"function_count"`
}

// AnalyzeComplexity calculates complexity metrics for a given file.
// Currently supports detailed AST analysis for Go, and line-based heuristics for others.
func AnalyzeComplexity(path string) (*ComplexityStats, error) {
	ext := filepath.Ext(path)
	if ext == ".go" {
		return analyzeGoFile(path)
	}
	return analyzeGenericFile(path)
}

func analyzeGoFile(path string) (*ComplexityStats, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse go file: %w", err)
	}

	stats := &ComplexityStats{}
	
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		stats.NodeCount++
		
		switch t := n.(type) {
		case *ast.FuncDecl:
			stats.Functions++
			stats.Cyclomatic += 1 // Base complexity for function
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause:
			stats.Cyclomatic++
		case *ast.BinaryExpr:
			if t.Op == token.LAND || t.Op == token.LOR {
				stats.Cyclomatic++
			}
		}
		return true
	})

	// Simple Halstead approximation: Volume = Length * log2(Vocabulary)
	// We'll use NodeCount as Length and distinct node types as Vocabulary approximation
	// For MVP, just scaling NodeCount is a proxy.
	stats.HalsteadVolume = float64(stats.NodeCount) * 5.0 // Placeholder multiplier

	return stats, nil
}

func analyzeGenericFile(path string) (*ComplexityStats, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(string(content), "\n")
	nonEmptyLines := 0
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			nonEmptyLines++
		}
	}

	// Heuristic for generic files
	return &ComplexityStats{
		NodeCount:      nonEmptyLines,
		Cyclomatic:     1 + (nonEmptyLines / 10), // Rough proxy
		HalsteadVolume: float64(len(content)),
		Functions:      0,
	}, nil
}
