package main

import (
	"bytes"
	"flag"
	"fmt"
	mutils "github.com/hertz-contrib/migrate/pkg/common/utils"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"github.com/hertz-contrib/migrate/pkg/common/mconsts"
	"github.com/hertz-contrib/migrate/pkg/visitor"

	"golang.org/x/tools/go/ast/astutil"
)

var (
	rootPath string
	mode     string
)

func init() {
	//mutils.FormatCode("/Users/skyenought/gopath/src/github.com/Skyenought/migrate")
	flag.StringVar(&rootPath, "root", "./testfile.go", "root path")
	flag.StringVar(&mode, "mode", "print", "switch tool mode")
}

func main() {
	var err error
	flag.Parse()

	fset := token.NewFileSet()
	path, _ := filepath.Abs(rootPath)
	file, _ := parser.ParseFile(fset, path, nil, parser.AllErrors)

	v := visitor.NewVisitor(fset, file)

	astutil.Apply(file, func(c *astutil.Cursor) bool {
		// gin -> server
		{
			v.RewriteImport(mconsts.GinPkg, mconsts.HertzServerPkg)
			v.ReplaceServerNew(c)
		}
		// rename gin.HandlerFunc -> app.HandlerFunc
		{
			// add context
			v.AddImport(mconsts.ContextPkg)
			v.AddImport(mconsts.HertzAppPkg)
			v.ChangeReqCtxSignatureInLine(c)
			v.ChangeReqCtxSignature(c)
		}
		// inner handlerFunc
		{
			v.AddImport(mconsts.HertzUtils)
			v.ReplaceUtilsH(c)
			v.ReplaceGinRequestMethod(c)
		}
		v.ReplaceHertzSpin(c)
		return true
	}, nil)

	// add server options
	astutil.Apply(file, func(c *astutil.Cursor) bool {
		v.AddServerOptions(c)
		return true
	}, nil)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		panic(err)
	}
	file, err = parser.ParseFile(fset, "", buf.String(), parser.ParseComments)
	if err != nil {
		panic(err)
	}

	switch mode {
	case "ast":
		ast.Print(fset, file)
	case "print":
		var output []byte
		buffer := bytes.NewBuffer(output)
		_ = format.Node(buffer, fset, file)
		fmt.Println(buffer.String())
	case "file":
		var output bytes.Buffer
		err = format.Node(&output, fset, file)
		if err != nil {
			log.Fatal(err)
		}
		// 写回原始文件
		err = os.WriteFile(path, output.Bytes(), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		// 输出Go代码
		fmt.Println("Cod e has been updated and saved to the original file.")
		mutils.FormatCode("/Users/skyenought/gopath/src/github.com/Skyenought/migrate")
	}
}
