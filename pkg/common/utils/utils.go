package mutils

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bytedance/gopkg/util/logger"
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
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})
	if err != nil {
		logger.Errorf("get golang file error: %s", err.Error())
	}
	return goFiles
}

func JudgeFuncParam(field *ast.Field, funcParam string) bool {
	starExpr, isStarExpr := field.Type.(*ast.StarExpr)
	selExpr, isSelExpr := starExpr.X.(*ast.SelectorExpr)
	tmp := strings.Split(funcParam, ".")
	packageName, StructName := tmp[0], tmp[1]
	return isStarExpr && isSelExpr && selExpr.X.(*ast.Ident).Name == packageName && selExpr.Sel.Name == StructName
}

func IsSrvRequestFunc(funcIdent *ast.Ident) bool {
	if funcIdent == nil || funcIdent.Name == "" {
		return false
	}

	return funcIdent.Name == "GET" ||
		funcIdent.Name == "POST" ||
		funcIdent.Name == "PUT" ||
		funcIdent.Name == "DELETE" ||
		funcIdent.Name == "PATCH" ||
		funcIdent.Name == "HEAD" ||
		funcIdent.Name == "OPTIONS" ||
		funcIdent.Name == "Any"
}

func GetImportPaths(file *ast.File) []string {
	var paths []string
	for _, importSpec := range file.Imports {
		// 获取导入路径
		importPath, err := strconv.Unquote(importSpec.Path.Value)
		if err != nil {
			logger.Errorf("get import Paths error: %s", err.Error())
			return nil
		}
		paths = append(paths, importPath)
	}
	return paths
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

func GetLastWord(s string) string {
	split := strings.Split(s, "/")
	lastWord := split[len(split)-1]
	return lastWord
}
