package main

import (
	"os"

	"github.com/dexterp/maker/internal"
	"github.com/dexterp/maker/internal/di"
	"github.com/dexterp/maker/internal/options"
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
