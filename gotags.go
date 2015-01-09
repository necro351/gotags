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
			handlePackageIdent(fset, x)
		case *ast.FuncDecl:
			handleFunctionIdent(fset, x)
			// Do not descend into function bodies
			return false
		case *ast.ValueSpec:
			// Do not descend into the value portion of vars
			for _, name := range x.Names {
				// Do not descend into function bodies
				handleObjectIdent(fset, name)
			}
			return false
		default:
			// But go everywhere else
			handleObjectIdent(fset, n)
		}
		return true
	})

	return nil
}

func handleFunctionIdent(fset *token.FileSet, decl *ast.FuncDecl) {
	position := fset.Position(decl.Type.Func)
	printIdent(decl.Name.Name, position.Filename, position.Line)
}

func handlePackageIdent(fset *token.FileSet, file *ast.File) {
	position := fset.Position(file.Package)
	printIdent(file.Name.Name, position.Filename, position.Line)
}

func handleObjectIdent(fset *token.FileSet, n interface{}) {
	switch x := n.(type) {
	case *ast.Ident:
		if x.Obj != nil && x.Obj.Decl != nil && expectedDecl(x.Obj.Decl) {
			pos := fset.Position(x.NamePos)
			printIdent(x.Name, pos.Filename, pos.Line)
		}
	}
}

func printIdent(name, file string, line int) {
	tags = append(tags, fmt.Sprintf("%s\t%s\t:%d", name, file, line))
}

func expectedDecl(decl interface{}) bool {
	switch decl.(type) {
	case *ast.TypeSpec, *ast.Field, *ast.FuncDecl, *ast.ValueSpec:
		return true
	}
	return false
}
