package visitor

import (
	"go/ast"
	"go/token"
)

var addr string

type Visitor struct {
	fset *token.FileSet
	f    *ast.File
}

func NewVisitor(fset *token.FileSet, f *ast.File) *Visitor {
	return &Visitor{fset: fset, f: f}
}
