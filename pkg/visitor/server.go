package visitor

import (
	mconsts "github.com/hertz-contrib/migrate/pkg/common/consts"
	"go/ast"

	"github.com/hertz-contrib/migrate/pkg/common/utils"

	"golang.org/x/tools/go/ast/astutil"
)

func (v *Visitor) ReplaceGinNew2ServerNew(c *astutil.Cursor) {
	// gin -> server
	call, ok := c.Node().(*ast.CallExpr)
	if ok {
		alias := v.FindImportAlias(mconsts.HertzServerPkg)
		identName := mutils.GetLastWord(mconsts.HertzServerPkg)
		if alias != "" {
			identName = alias
		}
		if mutils.IsPkgDot(call.Fun, "gin", "New") {
			c.Replace(&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent(identName),
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
			ident, ok := expr.X.(*ast.Ident)
			if !ok {
				return
			}
			alias := v.FindImportAlias(mconsts.HertzServerPkg)
			identName := mutils.GetLastWord(mconsts.HertzServerPkg)
			if alias != "" {
				identName = alias
			}
			if ident.Name == identName && expr.Sel.Name == "New" {
				withHostPorts := &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent(alias),
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
