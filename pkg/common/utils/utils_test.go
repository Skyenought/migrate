package mutils

import (
	"github.com/bytedance/gopkg/util/logger"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func InitAstFile() (*ast.File, *token.FileSet) {
	fset := token.NewFileSet()
	path, _ := filepath.Abs("./utils_test.go")
	file, _ := parser.ParseFile(fset, path, nil, parser.AllErrors)
	return file, fset
}

func TestGetAllGoFiles(t *testing.T) {
	files := GetAllGoFiles("./")
	assert.Equal(t, len(files), 3)
}

func TestGetImportPaths(t *testing.T) {
	file, _ := InitAstFile()
	paths := GetImportPaths(file)
	for _, path := range paths {
		logger.Info(path)
	}
	assert.Equal(t, 7, len(paths))
}
