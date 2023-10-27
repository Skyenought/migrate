package visitor

import (
	mutils "github.com/hertz-contrib/migrate/pkg/common/utils"
	"github.com/hertz-contrib/migrate/pkg/global"
	"go/ast"
	"math/rand"
	"time"

	"golang.org/x/tools/go/ast/astutil"
)

// DeleteImport 删除 import 语句
func (v *Visitor) DeleteImport(path string) {
	astutil.DeleteImport(v.fset, v.f, path)
}

// RewriteImport 重写 import 语句
func (v *Visitor) RewriteImport(oldPath, newPath string) {
	astutil.RewriteImport(v.fset, v.f, oldPath, newPath)
}

// AddImport is used to add an import path to a Go source file.
// If there is a package name conflict with the import path, it adds an alias to resolve the conflict.
// If no package name conflict exists, it adds the import path directly.
//
// Parameters:
//   - path: The import path of the package to be added.
//
// Note: This function modifies the provided AST file.
func (v *Visitor) AddImport(path string) {
	// Check if the import already exists
	if ImportExists(v.f, path) {
		return
	}
	var foundConflict = false
	getlastWorld := func(s string) string {
		return mutils.GetLastWord(s)
	}
	// Get all import paths in the file
	originPaths := mutils.GetImportPaths(v.f)
	for _, oPath := range originPaths {
		// Check if the package name conflicts with the import path
		if getlastWorld(oPath) == getlastWorld(path) {
			// If a conflict exists, generate an alias for the import path and add it
			astutil.AddNamedImport(v.fset, v.f, generateAlias(generateAlias(path)), path)
			foundConflict = true
			// Add the alias to the global alias map
			global.AliasMap[path] = generateAlias(generateAlias(path))
			break
		}
	}

	if !foundConflict {
		// If there is no package name conflict, add the import path directly
		astutil.AddImport(v.fset, v.f, path)
		return
	}
}

func generateAlias(s string) string {
	getlastWorld := func(s string) string {
		return mutils.GetLastWord(s)
	}
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)

	// 生成随机字母
	randomLetter := generateRandomLetter(generator)
	alias := string(randomLetter) + getlastWorld(s)
	return alias
}

func generateRandomLetter(generator *rand.Rand) rune {
	// Alphabet
	letters := "abcdefghijklmnopqrstuvwxyzdfjfsf"

	// Generate a random index
	randomIndex := generator.Intn(len(letters))

	// Get the random letter
	randomLetter := rune(letters[randomIndex])

	return randomLetter
}

func ImportExists(f *ast.File, path string) bool {
	for _, spec := range f.Imports {
		if spec.Path.Value == `"`+path+`"` {
			return true
		}
	}
	return false
}

// FindImportAlias finds the alias of the import path.
func (v *Visitor) FindImportAlias(importPath string) string {
	for _, imp := range v.f.Imports {
		if imp.Path.Value == `"`+importPath+`"` {
			if imp.Name != nil {
				return imp.Name.Name
			}
		}
	}
	return ""
}
