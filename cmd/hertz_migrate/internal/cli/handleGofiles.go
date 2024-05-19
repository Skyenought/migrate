package cli

import (
	"bytes"
	"go/ast"
	"go/parser"
	"log"
	"os"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal"
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/logic/chi"
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/logic/gin"
	nethttp "github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/logic/netHttp"
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/logs"
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/types"
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/utils"

	"golang.org/x/tools/go/ast/astutil"
)

// TODO:将 use-XXX 拆分为具体函数
func handleGoFiles(filepaths []string) error {
	for _, path := range filepaths {
		var containsGin bool

		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			logs.Debugf("Parse file fail, error: %v", err)
			return internal.ErrParseFile
		}

		astutil.AddNamedImport(fset, file, "hzserver", globalArgs.HzRepo+"/pkg/app/server")
		astutil.AddNamedImport(fset, file, "hzapp", globalArgs.HzRepo+"/pkg/app")
		astutil.AddNamedImport(fset, file, "hzroute", globalArgs.HzRepo+"/pkg/route")
		astutil.AddNamedImport(fset, file, "hzerrors", globalArgs.HzRepo+"/pkg/common/errors")
		astutil.AddNamedImport(fset, file, "hzutils", globalArgs.HzRepo+"/pkg/common/utils")

		for _, importSpec := range file.Imports {
			importStr := importSpec.Path.Value

			if globalArgs.UseGin {
				if strings.Contains(importStr, `github.com/gin-gonic/gin`) {
					containsGin = true
				}

				if strings.Contains(importStr, `"github.com/gin-contrib/cors"`) {
					importSpec.Path.Value = `"github.com/hertz-contrib/cors"`
					containsGin = true
				}

				if strings.Contains(importStr, `github.com/swaggo/gin-swagger`) {
					importSpec.Path.Value = `"github.com/hertz-contrib/swagger"`
					containsGin = true
				}
			}
		}

		if !containsGin && globalArgs.UseGin {
			continue
		}

		astutil.Apply(file, func(cursor *astutil.Cursor) bool {
			if globalArgs.UseGin {
				switch node := cursor.Node().(type) {
				case *ast.StarExpr:
					if sel, ok := node.X.(*ast.SelectorExpr); ok {
						if utils.CheckSelPkgAndStruct(sel, "gin", "Engine") {
							cursor.Replace(types.StarServerHertz)
						}

						if utils.CheckSelPkgAndStruct(sel, "gin", "RouterGroup") {
							cursor.Replace(types.StarRouteGroup)
						}
					}
				case *ast.FieldList:
					gin.ReplaceGinCtx(node)
				case *ast.SelectorExpr:
					if utils.CheckSelPkgAndStruct(node, "route", "IRoutes") {
						cursor.Replace(types.SelIRoutes)
					}
				}
			}

			switch node := cursor.Node().(type) {
			case *ast.StarExpr:
				if globalArgs.UseChi {
					if utils.CheckPtrPkgAndStructName(node, "chi", "Mux") {
						cursor.Replace(types.StarServerHertz)
					}
				}
			}

			if globalArgs.UseNetHTTP {
				nethttp.GetOptionsFromHttpServer(cursor, internal.GlobalHashMap)
				nethttp.PackServerHertz(cursor, internal.GlobalHashMap)
				nethttp.ReplaceNetHttpHandler(cursor)
			}

			return true
		}, nil)

		astutil.Apply(file, func(c *astutil.Cursor) bool {
			netHttpGroup(c, internal.WebCtxSet)

			switch node := c.Node().(type) {
			case *ast.SelectorExpr:
				if globalArgs.UseNetHTTP {
					if utils.CheckSelObj(node, "http", "ResponseWriter") {
						switch node.Sel.Name {
						case "WriteHeader":
							c.Replace(types.SelSetStatusCode)
						case "Write":
							c.Replace(types.SelWrite)
						case "Header":
							c.Replace(types.SelRespHeader)
						}
					}

					if node.Sel.Name == "HandleFunc" {
						node.Sel.Name = "Any"
					}
					nethttp.ReplaceRequestOp(node, c)
				}

				if globalArgs.UseGin {
					if utils.CheckSelPkgAndStruct(node, "gin", "HandlerFunc") {
						c.Replace(types.SelAppHandlerFunc)
					}

					if utils.CheckSelPkgAndStruct(node, "gin", "H") {
						node.X.(*ast.Ident).Name = "hzutils"
					}

					gin.ReplaceBinding(node, c)
					gin.ReplaceRequestOp(node, c)
					gin.ReplaceRespOp(node, c)
					gin.ReplaceErrorType(node)
				}
			case *ast.CallExpr:
				if globalArgs.UseChi {
					chi.PackChiRouterMethod(node)
					chi.PackChiNewRouter(node, c)
				}
				if globalArgs.UseNetHTTP {
					nethttp.ReplaceHttpOp(node, c)
					nethttp.ReplaceReqOrRespOp(node, c)
					nethttp.ReplaceReqURLQuery(node)
					if utils.CheckCallPkgAndMethodName(node, "http", "NotFound") {
						c.Replace(types.CallNotFound)
					}
				}

				if globalArgs.UseGin {
					gin.ReplaceGinNew(node, c)
					gin.ReplaceGinRun(node)
					gin.ReplaceGinCtxOp(node, c)
					gin.ReplaceCallReqOrResp(node, c)
					gin.ReplaceStatisFS(node)
				}
			}
			if globalArgs.UseGin {
				gin.ReplaceCtxParamList(c)
			}
			return true
		}, nil)

		var buf bytes.Buffer

		if err = printerConfig.Fprint(&buf, fset, file); err != nil {
			logs.Debugf("Fprint fail, error: %v", err)
			return internal.ErrSaveChanges
		}

		if err := os.WriteFile(path, buf.Bytes(), os.ModePerm); err == nil {
			log.Println("File updated:", path)
		}
	}
	return nil
}

func netHttpGroup(c *astutil.Cursor, funcSet mapset.Set[string]) {
	if globalArgs.UseNetHTTP {
		nethttp.PackFprintf(c)
		nethttp.ReplaceReqHeader(c)
		nethttp.ReplaceReqHeaderOperation(c)
		nethttp.ReplaceRespWrite(c)
		nethttp.ReplaceReqFormGet(c)
		nethttp.ReplaceReqFormValue(c)
		nethttp.ReplaceReqMultipartForm(c)
		nethttp.PackType2AppHandlerFunc(c)
		nethttp.ReplaceReqMultipartFormOperation(c, internal.GlobalHashMap)
		nethttp.ReplaceFuncBodyHttpHandlerParam(c, funcSet)
	}
}
