package visitor

import (
	mutils "github.com/hertz-contrib/migrate/pkg/common/utils"
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"
)

// ReplaceGinRequestMethod c.Request.Method -> string(c.Request.Method())
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

// ReplaceGinRequestFormValue c.Request.FormValue("") -> string(c.Request.FormValue(""))
func (v *Visitor) ReplaceGinRequestFormValue(c *astutil.Cursor) {
	expr, ok := c.Node().(*ast.CallExpr)
	if !ok {
		return
	}

	if selExpr, ok := expr.Fun.(*ast.SelectorExpr); ok {
		if innerExpr, ok := selExpr.X.(*ast.SelectorExpr); ok {
			// eg c Request FormValue "param"
			firstIdent := innerExpr.X.(*ast.Ident)
			secondIdent := innerExpr.Sel
			thirdIdent := selExpr.Sel

			if secondIdent.Name == "Request" && thirdIdent.Name == "FormValue" {
				c.Replace(&ast.CallExpr{
					Fun: &ast.Ident{Name: "string"},
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   firstIdent, // 将 innerExpr.X 改为 innerExpr
								Sel: thirdIdent,
							},
							Args: expr.Args,
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
		if isSelector {
			ident, isIdent := selectorExpr.X.(*ast.Ident)
			if isIdent {
				if selectorExpr.Sel.Name == "Next" {
					if ident.Name == "c" {
						ctxIdent := &ast.Ident{Name: "ctx"}
						callExpr.Args = []ast.Expr{ctxIdent}
					}
					if ident.Name == "ctx" {
						ctxIdent := &ast.Ident{Name: "c"}
						callExpr.Args = []ast.Expr{ctxIdent}
					} else {
						alias := string(mutils.GenerateRandomLetter()) + ident.Name
						ctxIdent := &ast.Ident{Name: alias}
						callExpr.Args = []ast.Expr{ctxIdent}
					}
				}
			}
		}
	}
}

func (v *Visitor) ReplaceGinShouldBindXxx(c *astutil.Cursor) {
	n := c.Node()
	if callExpr, ok := n.(*ast.CallExpr); ok {
		// 检查是否是 c.ShouldBindJSON(nil) 调用
		selectorExpr, isSelector := callExpr.Fun.(*ast.SelectorExpr)
		if isSelector {
			_, isIdent := selectorExpr.X.(*ast.Ident)
			if isIdent {
				switch selectorExpr.Sel.Name {
				case "ShouldBindJSON":
					selectorExpr.Sel.Name = "BindJSON"
				case "ShouldBindQuery":
					selectorExpr.Sel.Name = "BindQuery"
				case "ShouldBind":
					selectorExpr.Sel.Name = "Bind"
				case "ShouldBindHeader":
					selectorExpr.Sel.Name = "BindHeader"
				case "ShouldBindUri":
					selectorExpr.Sel.Name = "BindPath"
				case "ShouldBindYAML", "ShouldBindXML", "ShouldBindTOML":
					comment := &ast.Comment{
						Text:  "// TODO: unsupported this method",
						Slash: selectorExpr.Pos() - 1,
					}

					// 将注释添加到文件的注释列表
					v.f.Comments = append(v.f.Comments, &ast.CommentGroup{List: []*ast.Comment{comment}})
				}
			}
		}
	}
}
