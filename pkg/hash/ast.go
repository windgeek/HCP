package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"
)

// ComputeLogicHash calculates a hash based on the AST structure of a Go file.
// It ignores comments, whitespace, and import order.
func ComputeLogicHash(path string) (string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments) // Parse comments but we wont use them for hash
	if err != nil {
		return "", err
	}

	var sigBuilder strings.Builder

	// 1. Package Name
	fmt.Fprintf(&sigBuilder, "pkg:%s;", node.Name.Name)

	// 2. Imports (Sorted to ignore order)
	var imports []string
	for _, imp := range node.Imports {
		if imp.Path != nil {
			imports = append(imports, imp.Path.Value)
		}
	}
	sort.Strings(imports)
	for _, imp := range imports {
		fmt.Fprintf(&sigBuilder, "imp:%s;", imp)
	}

	// 3. Declarations (Funcs, Types, Vars)
	// We walk the AST and record "Logic Signatures"
	// For simplicity in v1, we hash:
	// - Function Names, Receiver, Params, Results
	// - Type Definitions
	// - Global Variables
	// We do NOT hash the body content deeply yet, just the structure signature to allow reformatting.
	// WAIT: The prompt says "AST-based hashing (ignoring whitespace, comments, and import orders)".
	// To be robust against reformatting inside function bodies, we should probably hash the *structure* of statements too,
	// but ignoring specific indentation/formatting.
	// Printing the AST node using printer.Fprint with standard config might be checking formatting too much.
	// Let's iterate top level decls and create a semantic signature.

	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			hashFunc(d, &sigBuilder)
		case *ast.GenDecl:
			hashGenDecl(d, &sigBuilder)
		}
	}

	// Hash the signature string
	h := sha256.New()
	h.Write([]byte(sigBuilder.String()))
	return hex.EncodeToString(h.Sum(nil)), nil
}

func hashFunc(fn *ast.FuncDecl, sb *strings.Builder) {
	// Receiver
	if fn.Recv != nil {
		sb.WriteString("recv:")
		for _, f := range fn.Recv.List {
			writeType(f.Type, sb)
		}
		sb.WriteString(";")
	}

	// Name
	fmt.Fprintf(sb, "func:%s;", fn.Name.Name)

	// Params
	sb.WriteString("params:")
	writeFieldList(fn.Type.Params, sb)
	sb.WriteString(";")

	// Results
	sb.WriteString("results:")
	writeFieldList(fn.Type.Results, sb)
	sb.WriteString(";")
	
	// For logic hash, we might want to capture Body roughly.
	// e.g. Count statements or hash their types?
	// For now, let's stick to Interface Hash (provenance of API).
	// If we want "Logic Preserved", we need body.
	// Let's walk the body non-recursively? No, full recursive is needed for "Logic".
	// Let's append a simplified representation of the body.
	if fn.Body != nil {
		sb.WriteString("body:{")
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			if n == nil {
				return true
			}
			// Just record the TYPE of node to capture control flow structure
			// e.g. IfStmt, AssignStmt, CallExpr
			// We ignore variable names inside body to allow renaming? 
			// No, renaming variables changes logic/intent usually.
			// Let's record the string representation of appropriate nodes but stripped of pos.
			switch t := n.(type) {
			case *ast.IfStmt:
				sb.WriteString("if;")
			case *ast.ForStmt:
				sb.WriteString("for;")
			case *ast.ReturnStmt:
				sb.WriteString("return;")
			case *ast.AssignStmt:
				sb.WriteString("assign;")
			case *ast.CallExpr:
				sb.WriteString("call;")
				// Hash function name being called
				if fun, ok := t.Fun.(*ast.Ident); ok {
					fmt.Fprintf(sb, "%s;", fun.Name)
				}
			}
			return true
		})
		sb.WriteString("};")
	}
}

func hashGenDecl(d *ast.GenDecl, sb *strings.Builder) {
	// Imports handled globally.
	if d.Tok == token.IMPORT {
		return
	}
	fmt.Fprintf(sb, "%s:", d.Tok)
	for _, spec := range d.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			fmt.Fprintf(sb, "type:%s;", s.Name.Name)
			writeType(s.Type, sb)
		case *ast.ValueSpec:
			for _, name := range s.Names {
				fmt.Fprintf(sb, "var:%s;", name.Name)
			}
			if s.Type != nil {
				writeType(s.Type, sb)
			}
		}
	}
	sb.WriteString(";")
}

func writeFieldList(fields *ast.FieldList, sb *strings.Builder) {
	if fields == nil {
		return
	}
	for _, f := range fields.List {
		writeType(f.Type, sb)
	}
}

func writeType(expr ast.Expr, sb *strings.Builder) {
	switch t := expr.(type) {
	case *ast.Ident:
		sb.WriteString(t.Name)
	case *ast.StarExpr:
		sb.WriteString("*")
		writeType(t.X, sb)
	case *ast.SelectorExpr:
		writeType(t.X, sb)
		sb.WriteString(".")
		sb.WriteString(t.Sel.Name)
	case *ast.ArrayType:
		sb.WriteString("[]")
		writeType(t.Elt, sb)
	// Add more complex types as needed
	default:
		sb.WriteString("T")
	}
	sb.WriteString(",")
}
