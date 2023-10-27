package visitor

import (
	mconsts "github.com/hertz-contrib/migrate/pkg/common/consts"
	"go/ast"

	"github.com/hertz-contrib/migrate/pkg/common/utils"

	"golang.org/x/tools/go/ast/astutil"
)

func (v *Visitor) ReplaceGinH2UtilsH(c *astutil.Cursor) {
	sel, ok := c.Node().(*ast.SelectorExpr)
	if ok {
		// 如果是 gin.H 的形式，替换为 utils.H
		if mutils.IsPkgDot(sel, "gin", "H") {
			pkgName := "utils"
			// find import alias
			alias := v.FindImportAlias(mconsts.HertzUtils)
			ident := ast.NewIdent(pkgName)
			if alias != "" {
				pkgName = alias
			}
			c.Replace(&ast.SelectorExpr{
				X:   ident,
				Sel: ast.NewIdent("H"),
			})
		}
	}
}
