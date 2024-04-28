package cli

import (
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/logs"
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/utils"
	"github.com/urfave/cli/v2"
	"go/token"
)

func preload(ctx *cli.Context) (files []string, gomods []string, err error) {
	var (
		gofiles   []string
		goModDirs []string
	)
	fset = token.NewFileSet()
	globalArgs.IgnoreDirs = ctx.StringSlice("ignore-dirs")

	if globalArgs.UseChi {
		globalArgs.UseNetHTTP = true
	}

	if globalArgs.Debug {
		logs.SetLevel(logs.LevelDebug)
	}

	if globalArgs.TargetDir != "" {
		gofiles, err = utils.CollectGoFiles(globalArgs.TargetDir, globalArgs.IgnoreDirs)
		if err != nil {
			return
		}

		gomods, err = utils.SearchAllDirHasGoMod(globalArgs.TargetDir)
		if err != nil {
			return
		}

		for _, dir := range goModDirs {
			wg.Add(1)
			dir := dir
			go func() {
				defer wg.Done()
				utils.RunGoGet(dir, globalArgs.HzRepo)
			}()
		}
		wg.Wait()
	}
	return gofiles, gomods, err
}
