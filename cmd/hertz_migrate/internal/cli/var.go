package cli

import (
	"go/printer"
	"go/token"
	"sync"
)

var (
	globalArgs = &Args{}
	fset       *token.FileSet
	wg         sync.WaitGroup

	printerConfig = printer.Config{
		Mode:     printer.UseSpaces,
		Tabwidth: 4,
	}
)
