package cli

import (
	"go/parser"

	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal"
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/logic"
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/logic/gin"
	nethttp "github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/logic/netHttp"
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/logs"
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/utils"
	"github.com/urfave/cli/v2"

	"golang.org/x/tools/go/ast/astutil"
)

func Run(c *cli.Context) error {
	var (
		gofiles   []string
		goModDirs []string
		err       error
	)

	if gofiles, goModDirs, err = preload(c); err != nil {
		return err
	}

	for _, path := range gofiles {
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			logs.Debugf("Parse file fail, error: %v", err)
			return internal.ErrParseFile
		}

		// collect global information
		astutil.Apply(file, func(c *astutil.Cursor) bool {
			logic.GetHttpServerProps(c)
			if globalArgs.UseGin {
				gin.GetFuncNameHasGinCtx(c)
			}
			if globalArgs.UseNetHTTP {
				nethttp.FindHandlerFuncName(c, internal.WebCtxSet)
			}
			return true
		}, nil)
	}

	if err = handleGoFiles(gofiles); err != nil {
		return err
	}

	for _, dir := range goModDirs {
		utils.RunGoImports(dir)
		utils.RunGoModTidy(dir)
	}
	logs.Info("everything are ok!")
	return nil
}
