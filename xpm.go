package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	xpm "xpm_internal"
)

var _Version string

var (
	hflag  bool
	vflag  bool
	prefix string
)

func initGlobal() {
	_Version = "0.0.1"
}

func initFlags() {
	flag.BoolVar(&hflag, "help", false, "show this help information")
	flag.BoolVar(&hflag, "h", false, "")
	flag.BoolVar(&vflag, "version", false, "show version and exit")
	flag.BoolVar(&vflag, "v", false, "")
	flag.StringVar(&prefix, "prefix", ".", "set `prefix` path for installation")
	flag.StringVar(&prefix, "p", ".", "")

	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, `XqPackageManager version: xpm/%s
Usage: xpm [-hv] [-p prefix] install <package_name>[@version] [config]

Options:
`, _Version)
	flag.PrintDefaults()
}

func main() {
	initGlobal()
	initFlags()

	if hflag {
		usage()
		return
	}

	if vflag {
		fmt.Fprintf(os.Stdout, "xpm version %s\n", _Version)
		return
	}

	args := flag.Args()
	if len(args) < 2 {
		usage()
		return
	}
	switch args[0] {
	case "install":
		var vers string
		vlist := strings.SplitN(args[1], "@", 2)
		if len(vlist) == 2 {
			vers = vlist[1]
		} else {
			vers = "latest"
		}
		config := args[2:]
		err := xpm.Install(vlist[0], vers, config, prefix)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
	default:
		usage()
		return
	}
}
