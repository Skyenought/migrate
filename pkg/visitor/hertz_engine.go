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
