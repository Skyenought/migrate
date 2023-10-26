package visitor

import (
	"go/ast"

	"github.com/hertz-contrib/migrate/pkg/common/utils"

	"golang.org/x/tools/go/ast/astutil"
)

func (v *Visitor) ReplaceServerNew(c *astutil.Cursor) {
	// gin -> server
	call, ok := c.Node().(*ast.CallExpr)
	if ok {
		if mutils.IsPkgDot(call.Fun, "gin", "New") {
			c.Replace(&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("server"),
					Sel: ast.NewIdent("New"),
				},
			})
		}
	}
}

func (v *Visitor) AddServerOptions(c *astutil.Cursor) {
	v.AddWithHostPorts(c)
}

func (v *Visitor) AddWithHostPorts(c *astutil.Cursor) {
	if callExpr, ok := c.Node().(*ast.CallExpr); ok {
		if expr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if expr.Sel.Name == "New" && expr.X.(*ast.Ident).Name == "server" {
				withHostPorts := &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("server"),
						Sel: ast.NewIdent("WithHostPorts"),
					},
					Args: []ast.Expr{ast.NewIdent(addr)},
				}

				// 插入 withHostPorts 到参数列表
				callExpr.Args = append(callExpr.Args, withHostPorts)
			}
		}
	}
}
