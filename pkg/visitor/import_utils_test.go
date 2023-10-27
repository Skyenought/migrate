package visitor

import (
	"bytes"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"go/format"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
)

func TestAddImport(t *testing.T) {
	fset := token.NewFileSet()
	path, _ := filepath.Abs("./import_utils_test.go")
	file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		t.Fatal(err)
	}
	_ = utils.H{}
	v := NewVisitor(fset, file)
	v.AddImport("github.com/hertz-contrib/migrate/pkg/common/utils")

	var output []byte
	buffer := bytes.NewBuffer(output)
	err = format.Node(buffer, fset, file)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(buffer.String())
}
