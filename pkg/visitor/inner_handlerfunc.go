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
func (v *Visitor) ReplaceGinNext(c *astutil.Cursor) {
	n := c.Node()
	if callExpr, ok := n.(*ast.CallExpr); ok {
		// 检查是否是 c.Next() 调用
		selectorExpr, isSelector := callExpr.Fun.(*ast.SelectorExpr)
		ident, isIdent := selectorExpr.X.(*ast.Ident)

		if isSelector && isIdent {
			if selectorExpr.Sel.Name == "Next" {
				if ident.Name == "c" {
					ctxIdent := &ast.Ident{Name: "ctx"}
					callExpr.Args = []ast.Expr{ctxIdent}
				}
				if ident.Name == "ctx" {
					ctxIdent := &ast.Ident{Name: "c"}
					callExpr.Args = []ast.Expr{ctxIdent}
				}
			}
		}
	}
}
