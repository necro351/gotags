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

	// Parse the file containing this very example
	// but stop after processing the imports.
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return err
	}

	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			if x.Obj != nil && x.Obj.Decl != nil {
				pos := fset.Position(x.NamePos)
				tags = append(tags, fmt.Sprintf("%s\t%s\t:%d", x.Name, pos.Filename, pos.Line))
			}
		}
		return true
	})

	return nil
}
