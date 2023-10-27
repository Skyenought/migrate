package visitor

import (
	"go/ast"

	mconsts "github.com/hertz-contrib/migrate/pkg/common/consts"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/hertz-contrib/migrate/pkg/common/utils"
)

// ChangeReqCtxSignature 修改请求上下文的签名
// func (c *gin.Context) -> func(ctx context.Context, c *app.RequestContext)
func (v *Visitor) ChangeReqCtxSignature(c *astutil.Cursor) {
	funcDecl, ok := c.Node().(*ast.FuncDecl)
	if !ok {
		return
	}
	for _, field := range funcDecl.Type.Params.List {
		if len(field.Names) == 1 {
			if mutils.JudgeFuncParam(field, mconsts.GinCtx) {
				// 检查 *gin.Context 参数
				switch field.Names[0].Name {
				case "c":
					v.ReplaceHandlerFuncParams(funcDecl, "ctx", "c")
				case "ctx":
					v.ReplaceHandlerFuncParams(funcDecl, "c", "ctx")
				default:
					v.ReplaceHandlerFuncParams(funcDecl, "ctx", field.Names[0].Name)
				}
			}
		}
	}
}

// ChangeReqCtxSignatureInLine 修改请求上下文的签名
func (v *Visitor) ChangeReqCtxSignatureInLine(c *astutil.Cursor) {
	/*
		遇到 h.GET("/ping", func(c *gin.Context) {}) 这种情况
	*/
	exprStmt, ok := c.Node().(*ast.ExprStmt)
	if ok {
		callExpr, ok := exprStmt.X.(*ast.CallExpr)
		if ok {
			args := callExpr.Args
			if args == nil {
				return
			}
			// remove route path
			args = args[1:]
			for _, arg := range args {
				funcLit, ok := arg.(*ast.FuncLit)
				if ok {
					expr := funcLit.Type.Params.List[0]
					paramName := expr.Names[0].Name
					if mutils.JudgeFuncParam(expr, mconsts.GinCtx) {
						switch paramName {
						case "c":
							v.ReplaceHandlerFuncParamsByLit(funcLit, "ctx", "c")
						case "ctx":
							v.ReplaceHandlerFuncParamsByLit(funcLit, "c", "ctx")
						default:
							v.ReplaceHandlerFuncParamsByLit(funcLit, "ctx", expr.Names[0].Name)
						}
					}
				}
			}
		}
	}

	/*
		遇到 return func(c context.Context, ctx *app.RequestContext) 这种情况
	*/
	funcLit, ok := c.Node().(*ast.FuncLit)
	if ok {
		for _, field := range funcLit.Type.Params.List {
			starExpr, ok := field.Type.(*ast.StarExpr)
			if ok {
				selExpr, ok := starExpr.X.(*ast.SelectorExpr)
				if ok {
					if selExpr.X.(*ast.Ident).Name == "gin" && selExpr.Sel.Name == "Context" {
						paramName := field.Names[0].Name
						if mutils.JudgeFuncParam(field, mconsts.GinCtx) {
							switch paramName {
							case "c":
								v.ReplaceHandlerFuncParamsByLit(funcLit, "ctx", "c")
							case "ctx":
								v.ReplaceHandlerFuncParamsByLit(funcLit, "c", "ctx")
							default:
								v.ReplaceHandlerFuncParamsByLit(funcLit, "ctx", field.Names[0].Name)
							}
						}
					}
				}
			}
		}
	}
}

func (v *Visitor) ReplaceHandlerFuncParams(funcDecl *ast.FuncDecl, preCtx, newCtx string) {
	// create new param ctx context.Context
	newParam1 := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(preCtx)},
		Type:  &ast.Ident{Name: mconsts.NormalCtx},
	}

	// create new param c *app.RequestContext
	getlastWorld := func(s string) string {
		return mutils.GetLastWord(s)
	}

	pkgName := getlastWorld(mconsts.HertzAppPkg)
	// find import alias
	alias := v.FindImportAlias(mconsts.HertzAppPkg)
	ident := ast.NewIdent(pkgName)
	if alias != "" {
		ident = ast.NewIdent(alias)
	}

	newParam2 := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(newCtx)},
		Type: &ast.StarExpr{
			X: &ast.SelectorExpr{
				X:   ident,
				Sel: ast.NewIdent("RequestContext"),
			},
		},
	}

	// 替换参数列表
	funcDecl.Type.Params.List = []*ast.Field{newParam1, newParam2}
}

func (v *Visitor) ReplaceHandlerFuncParamsByLit(funcDecl *ast.FuncLit, preCtx, newCtx string) {
	// 创建新参数 ctx context.Context
	newParam1 := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(preCtx)},
		Type:  &ast.Ident{Name: mconsts.NormalCtx},
	}

	// create new param c *app.RequestContext
	getlastWorld := func(s string) string {
		return mutils.GetLastWord(s)
	}

	pkgName := getlastWorld(mconsts.HertzAppPkg)
	// find import alias
	alias := v.FindImportAlias(mconsts.HertzAppPkg)
	ident := ast.NewIdent(pkgName)
	if alias != "" {
		ident = ast.NewIdent(alias)
	}

	newParam2 := &ast.Field{
		Names:   []*ast.Ident{ast.NewIdent(newCtx)},
		Type:    &ast.StarExpr{X: &ast.SelectorExpr{X: ident, Sel: ast.NewIdent("RequestContext")}},
		Comment: nil,
	}

	// 替换参数列表
	funcDecl.Type.Params.List = []*ast.Field{newParam1, newParam2}
}
