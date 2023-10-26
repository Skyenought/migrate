package visitor

import (
	"go/ast"
	"golang.org/x/tools/go/ast/astutil"
)

func (v *Visitor) ReplaceGinRequestMethod(c *astutil.Cursor) {
	if selExpr, ok := c.Node().(*ast.SelectorExpr); ok {
		if innerExpr, ok := selExpr.X.(*ast.SelectorExpr); ok {
			firstIdent := innerExpr.X.(*ast.Ident)
			secondIdent := innerExpr.Sel
			thirdIdent := selExpr.Sel

			if secondIdent.Name == "Request" && thirdIdent.Name == "Method" {
				c.Replace(&ast.CallExpr{
					Fun: &ast.Ident{Name: "string"},
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   firstIdent,
									Sel: secondIdent,
								},
								Sel: thirdIdent,
							},
							Args: nil,
						},
					},
				})
			}
		}
	}
}
