// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logic

import (
	. "go/ast"

	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal"
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/types"
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/utils"
	"golang.org/x/tools/go/ast/astutil"
)

func collectHertzConfigOptions(elt Expr) {
	if kve, ok := elt.(*KeyValueExpr); ok {
		if ident, ok := kve.Key.(*Ident); ok {
			internal.HertzConfigOptions = append(internal.HertzConfigOptions, types.ExportServerOption("WithDisablePrintRoute", []Expr{NewIdent("true")}))
			switch ident.Name {
			case "Addr":
				internal.HertzConfigOptions = append(internal.HertzConfigOptions, types.ExportServerOption("WithHostPorts", []Expr{kve.Value}))
			case "WriteTimeout":
				internal.HertzConfigOptions = append(internal.HertzConfigOptions, types.ExportServerOption("WithWriteTimeout", []Expr{kve.Value}))
			case "ReadTimeout":
				internal.HertzConfigOptions = append(internal.HertzConfigOptions, types.ExportServerOption("WithReadTimeout", []Expr{kve.Value}))
			case "IdleTimeout":
				internal.HertzConfigOptions = append(internal.HertzConfigOptions, types.ExportServerOption("WithIdleTimeout", []Expr{kve.Value}))
			case "TLSConfig":
				internal.HertzConfigOptions = append(internal.HertzConfigOptions, types.ExportServerOption("WithTLS", []Expr{kve.Value}))
			}
		}
	}
}

func GetHttpServerProps(cur *astutil.Cursor) {
	if cLit, ok := cur.Node().(*CompositeLit); ok {
		if sel, ok := cLit.Type.(*SelectorExpr); ok {
			if utils.CheckSelPkgAndStruct(sel, "http", "Server") {
				if len(cLit.Elts) > 0 {
					for _, elt := range cLit.Elts {
						collectHertzConfigOptions(elt)
					}
				}
			}
		}
	}

	if callExpr, ok := cur.Node().(*CallExpr); ok {
		if sel, ok := callExpr.Fun.(*SelectorExpr); ok {
			if ident, ok := sel.X.(*Ident); ok {
				if utils.CheckObjSelExpr(ident.Obj, "gin", "New") ||
					utils.CheckObjSelExpr(ident.Obj, "gin", "Default") ||
					utils.CheckObjStarExpr(ident.Obj, "gin", "Engine") {
					if sel.Sel.Name == "Run" {
						if len(callExpr.Args) == 1 {
							internal.ServerPort = callExpr.Args[0]
						}
					}
				}
			}
		}
	}
}
