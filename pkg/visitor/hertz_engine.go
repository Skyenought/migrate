package visitor

import (
	"go/ast"

	"github.com/hertz-contrib/migrate/pkg/common/utils"

	"golang.org/x/tools/go/ast/astutil"
)

func (v *Visitor) ReplaceGinRun2HertzSpin(c *astutil.Cursor) {
	call, ok := c.Node().(*ast.CallExpr)
	if ok {
		if mutils.IsDot(call.Fun, "Run") && len(call.Args) == 1 {
			if ident, ok := call.Args[0].(*ast.Ident); ok {
				addr = ident.Name
			}
			c.Replace(&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   call.Fun.(*ast.SelectorExpr).X,
					Sel: ast.NewIdent("Spin"),
				},
				Args: nil,
			})
		}
	}
}

func (v *Visitor) ReplaceGinHandlerFunc(c *astutil.Cursor) {
	indent, ok := c.Node().(*ast.Ident)
	if !ok || indent.Obj == nil {
		return
	}

	funcDecl, ok := indent.Obj.Decl.(*ast.FuncDecl)
	if ok {
		if funcDecl.Type.Results == nil {
			return
		}
		for _, returnValue := range funcDecl.Type.Results.List {
			if selExpr, ok := returnValue.Type.(*ast.SelectorExpr); ok {
				if selExpr.X.(*ast.Ident).Name == "gin" && selExpr.Sel.Name == "HandlerFunc" {
					newSelExpr := &ast.SelectorExpr{
						X:   &ast.Ident{Name: "app"},
						Sel: &ast.Ident{Name: "HandlerFunc"},
					}
					returnValue.Type = newSelExpr
				}
			}
		}
	}
}
