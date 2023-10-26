package mutils

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/hertz-contrib/migrate/pkg/common/mconsts"
)

func DelElementFromSlice[T comparable](a []T, ele T) []T {
	j := 0
	for _, v := range a {
		if v != ele {
			a[j] = v
			j++
		}
	}
	return a[:j]
}

func GetAllGoFiles(dir string) []string {
	var goFiles []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})
	return goFiles
}

func JudgeFuncParam(field *ast.Field, funcParam string) bool {
	starExpr, isStarExpr := field.Type.(*ast.StarExpr)
	selExpr, isSelExpr := starExpr.X.(*ast.SelectorExpr)
	tmp := strings.Split(funcParam, ".")
	packageName, StructName := tmp[0], tmp[1]
	return isStarExpr && isSelExpr && selExpr.X.(*ast.Ident).Name == packageName && selExpr.Sel.Name == StructName
}

func ReplaceHandlerFuncParams(funcDecl *ast.FuncDecl, preCtx, newCtx string) {
	// 创建新参数 ctx context.Context
	newParam1 := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(preCtx)},
		Type:  &ast.Ident{Name: mconsts.NormalCtx},
	}

	// 创建新参数 c *app.RequestContext
	newParam2 := &ast.Field{
		Names:   []*ast.Ident{ast.NewIdent(newCtx)},
		Type:    &ast.StarExpr{X: &ast.SelectorExpr{X: ast.NewIdent("app"), Sel: ast.NewIdent("RequestContext")}},
		Comment: nil,
	}

	// 替换参数列表
	funcDecl.Type.Params.List = []*ast.Field{newParam1, newParam2}
}

func ReplaceHandlerFuncParamsByLit(funcDecl *ast.FuncLit, preCtx, newCtx string) {
	// 创建新参数 ctx context.Context
	newParam1 := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(preCtx)},
		Type:  &ast.Ident{Name: mconsts.NormalCtx},
	}

	// 创建新参数 c *app.RequestContext
	newParam2 := &ast.Field{
		Names:   []*ast.Ident{ast.NewIdent(newCtx)},
		Type:    &ast.StarExpr{X: &ast.SelectorExpr{X: ast.NewIdent("app"), Sel: ast.NewIdent("RequestContext")}},
		Comment: nil,
	}

	// 替换参数列表
	funcDecl.Type.Params.List = []*ast.Field{newParam1, newParam2}
}

func IsSrvRequestFunc(funcIdent *ast.Ident) bool {
	return funcIdent.Name == "GET" ||
		funcIdent.Name == "POST" ||
		funcIdent.Name == "PUT" ||
		funcIdent.Name == "DELETE" ||
		funcIdent.Name == "PATCH" ||
		funcIdent.Name == "HEAD" ||
		funcIdent.Name == "OPTIONS" ||
		funcIdent.Name == "Any"
}

func IsPkgDot(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && IsIdent(sel.X, pkg) && IsIdent(sel.Sel, name)
}

func IsIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}

func IsString(expr ast.Expr, value string) bool {
	lit, ok := expr.(*ast.BasicLit)
	return ok && lit.Kind == token.STRING && lit.Value == fmt.Sprintf("%q", value)
}

func IsDot(expr ast.Expr, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && IsIdent(sel.Sel, name)
}
