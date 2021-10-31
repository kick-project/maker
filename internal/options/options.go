package options

import (
	"github.com/docopt/docopt-go"
)

func GetUsage(argv []string, version string) *Options {
	usage := `maker

Usage:
  maker [--dotenv=<files>] [--scan=<makefile>] [--prefix=<prefix>] <target>
  maker --menu=<makefile>

Options:
  -h --help          Show this screen
  --version          Show version
  --dotenv=<files>   List of comma separated paths to load dotenv files [default: ~/.env,.env]
  --menu=<makefile>  Print Makefile menu
  --prefix=<prefix>  Scan for "<target>" then "<prefix><target>" when scanning makefiles [default: "_"]
  --scan=<makefile>  Scan make makefile for targets and wrap with dotenv variables [default: Makefile]
  <target>           makefile target
`

	opts, err := docopt.ParseArgs(usage, argv, version)
	if err != nil {
		panic(err)
	}
	config := &Options{}
	err = opts.Bind(config)
	if err != nil {
		panic(err)
	}
	return config
}

type Options struct {
	Args   []string `docopt:"<args>"`
	Cmd    string   `docopt:"<cmd>"`
	Dotenv string   `docopt:"--dotenv"`
	Menu   string   `docopt:"--menu"`
	Prefix string   `docopt:"--prefix"`
	Scan   string   `docopt:"--scan"`
	Target string   `docopt:"<target>"`
}
