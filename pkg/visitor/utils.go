package visitor

import "golang.org/x/tools/go/ast/astutil"

// DeleteImport 删除 import 语句
func (v *Visitor) DeleteImport(path string) {
	astutil.DeleteImport(v.fset, v.f, path)
}

// RewriteImport 重写 import 语句
func (v *Visitor) RewriteImport(oldPath, newPath string) {
	astutil.RewriteImport(v.fset, v.f, oldPath, newPath)
}

// AddImport 添加 import 语句
func (v *Visitor) AddImport(path string) {
	astutil.AddImport(v.fset, v.f, path)
}
