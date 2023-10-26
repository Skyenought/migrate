package visitor

import (
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/hertz-contrib/migrate/pkg/common/mconsts"
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
			// 检查 *gin.Context 参数
			switch field.Names[0].Name {
			case "c":
				if mutils.JudgeFuncParam(field, mconsts.GinCtx) {
					mutils.ReplaceHandlerFuncParams(funcDecl, "ctx", "c")
				}
			case "ctx":
				if mutils.JudgeFuncParam(field, mconsts.GinCtx) {
					mutils.ReplaceHandlerFuncParams(funcDecl, "c", "ctx")
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
		args := exprStmt.X.(*ast.CallExpr).Args
		// remove route path
		args = args[1:]
		for _, arg := range args {
			funcLit, ok := arg.(*ast.FuncLit)
			if ok {
				expr := funcLit.Type.Params.List[0]
				paramName := expr.Names[0].Name
				switch paramName {
				case "c":
					if mutils.JudgeFuncParam(expr, mconsts.GinCtx) {
						mutils.ReplaceHandlerFuncParamsByLit(funcLit, "ctx", "c")
					}
				case "ctx":
					if mutils.JudgeFuncParam(expr, mconsts.GinCtx) {
						mutils.ReplaceHandlerFuncParamsByLit(funcLit, "c", "ctx")
					}
				}
			}
		}
	}
}
