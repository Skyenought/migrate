package visitor

import (
	"go/ast"

	"github.com/hertz-contrib/migrate/pkg/common/utils"

	"golang.org/x/tools/go/ast/astutil"
)

func (v *Visitor) ReplaceUtilsH(c *astutil.Cursor) {
	sel, ok := c.Node().(*ast.SelectorExpr)
	if ok {
		// 如果是 gin.H 的形式，替换为 utils.H
		if mutils.IsPkgDot(sel, "gin", "H") {
			c.Replace(&ast.SelectorExpr{
				X:   ast.NewIdent("utils"),
				Sel: ast.NewIdent("H"),
			})
		}
	}
}
