package main

import (
	"os"

	"github.com/kick-project/maker/internal"
	"github.com/kick-project/maker/internal/di"
	"github.com/kick-project/maker/internal/options"
)

func main() {
	opts := options.GetUsage(os.Args[1:], internal.Version)
	inject := di.Defaults(&di.DI{Prefix: "_"})
	if len(opts.Menu) > 0 {
		menu := inject.MakeMenu()
		menu.Display(opts.Menu)
		return
	}

	denv := inject.MakeDotenv()
	denv.WrapTarget(opts.Dotenv, opts.Scan, opts.Target)
}
