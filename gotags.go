package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
)

var tags []string

func main() {
	err := filepath.Walk("./", parseFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	sort.Strings(tags)
	for _, v := range tags {
		fmt.Printf("%s\n", v)
	}
}

func parseFile(filePath string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	matches, err := filepath.Match("*.go", info.Name())
	if err != nil {
		return err
	}
	if !matches {
		return nil
	}

	fset := token.NewFileSet() // positions are relative to fset

	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil && f == nil {
		return nil
	}

	// Inspect the AST and print all identifier declarations as tags
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			handleIdent(fset, x.Name, false)
		case *ast.FuncDecl:
			// Do not descend into function bodies
			handleIdent(fset, x.Name, true)
			return false
		case *ast.ValueSpec:
			// Do not descend into the value portion of vars
			for _, name := range x.Names {
				// Do not descend into function bodies
				handleIdent(fset, name, true)
			}
			return false
		default:
			// But go everywhere else
			handleIdent(fset, n, true)
		}
		return true
	})

	return nil
}

func handleIdent(fset *token.FileSet, n interface{}, reqdecl bool) {
	switch x := n.(type) {
	case *ast.Ident:
		if !reqdecl || (x.Obj != nil && x.Obj.Decl != nil && expectedDecl(x.Obj.Decl)) {
			pos := fset.Position(x.NamePos)
			tags = append(tags, fmt.Sprintf("%s\t%s\t:%d", x.Name, pos.Filename, pos.Line))
		}
	}
}

func expectedDecl(decl interface{}) bool {
	switch decl.(type) {
	case *ast.TypeSpec, *ast.Field, *ast.FuncDecl, *ast.ValueSpec:
		return true
	}
	return false
}
