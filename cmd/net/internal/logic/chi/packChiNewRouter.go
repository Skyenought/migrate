package chi

import (
	"github.com/hertz-contrib/migrate/cmd/net/internal/global"
	nethttp "github.com/hertz-contrib/migrate/cmd/net/internal/logic/netHttp"
	. "go/ast"
	"golang.org/x/tools/go/ast/astutil"
)

func PackChiNewRouter(cur *astutil.Cursor) {
	stmt, ok := cur.Node().(*AssignStmt)
	if !ok || len(stmt.Lhs) != 1 || len(stmt.Rhs) != 1 {
		return
	}
	callExpr, ok := stmt.Rhs[0].(*CallExpr)
	if !ok {
		return
	}
	selExpr, ok := callExpr.Fun.(*SelectorExpr)
	if !ok {
		return
	}
	if selExpr.Sel.Name == "NewRouter" {
		callExpr.Fun = &SelectorExpr{
			X:   NewIdent("server"),
			Sel: NewIdent("Default"),
		}
		global.Map["server"] = stmt.Lhs[0].(*Ident).Name
		nethttp.AddOptionsForServer(callExpr, global.Map)
	}
}
