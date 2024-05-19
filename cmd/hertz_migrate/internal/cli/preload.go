package cli

import (
	"go/token"
	"os"

	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/logs"
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal/utils"
	"github.com/urfave/cli/v2"
)

func preload(ctx *cli.Context) (files, gomodDirs []string, err error) {
	fset = token.NewFileSet()
	globalArgs.IgnoreDirs = ctx.StringSlice("ignore-dirs")

	if globalArgs.UseChi {
		globalArgs.UseNetHTTP = true
	}

	if globalArgs.Debug {
		logs.SetLevel(logs.LevelDebug)
	}

	if globalArgs.TargetDir != "" {
		files, err = utils.CollectGoFiles(globalArgs.TargetDir, globalArgs.IgnoreDirs)
		if err != nil {
			return
		}

		gomodDirs, err = utils.SearchAllDirHasGoMod(globalArgs.TargetDir)
		if err != nil {
			return
		}

		for _, dir := range gomodDirs {
			wg.Add(1)
			dir := dir
			go func() {
				defer wg.Done()
				if err := utils.RunGoGet(dir, globalArgs.HzRepo); err != nil {
					logs.Errorf("go get hertz fail, %v", err)
					os.Exit(1)
				}
			}()
		}
		wg.Wait()
	}
	return files, gomodDirs, err
}
