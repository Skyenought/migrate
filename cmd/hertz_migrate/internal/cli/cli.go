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

package cli

import (
	"github.com/hertz-contrib/migrate/cmd/hertz_migrate/internal"
	"github.com/urfave/cli/v2"
)

type Args struct {
	TargetDir  string
	Filepath   string
	HzRepo     string
	IgnoreDirs []string
	Debug      bool
	UseGin     bool
	UseNetHTTP bool
	UseChi     bool
}

const ignoreDirsText = `
Fill in the folders to be ignored, separating the folders with ",".
Example:
    hertz_migrate -target-dir ./project -ignore-dirs=hz_gen -ignore-dirs=vendor
`

func Init() *cli.App {
	app := cli.NewApp()
	app.Name = "hertz_migrate"
	app.Usage = "A tool for migrating to hertz from other go web frameworks"
	app.Version = internal.Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "hz-repo",
			Aliases:     []string{"r"},
			Value:       "github.com/cloudwego/hertz",
			Usage:       "Specify the url of the hertz repository you want to bring in.",
			Destination: &globalArgs.HzRepo,
		},
		&cli.StringFlag{
			Name:        "target-dir",
			Aliases:     []string{"d"},
			Usage:       "Project directory you wants to migrate.",
			Destination: &globalArgs.TargetDir,
		},
		&cli.StringSliceFlag{
			Name:    "ignore-dirs",
			Aliases: []string{"I"},
			Usage:   ignoreDirsText,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"D"},
			Destination: &globalArgs.Debug,
			Value:       false,
		},
		&cli.BoolFlag{
			Name:        "use-gin",
			Aliases:     []string{"g"},
			Usage:       "Use this flag to migrate gin to the hertz framework.",
			Destination: &globalArgs.UseGin,
		},
		&cli.BoolFlag{
			Name:        "use-net-http",
			Aliases:     []string{"n"},
			Usage:       "Use this flag to migrate net/http to the hertz framework.",
			Destination: &globalArgs.UseNetHTTP,
		},
		&cli.BoolFlag{
			Name:        "use-chi",
			Aliases:     []string{"c"},
			Usage:       "Use this flag to migrate chi to the hertz framework.\nEnabling it will automatically enable the use-net-http flag",
			Destination: &globalArgs.UseNetHTTP,
		},
	}
	app.Action = Run
	return app
}
