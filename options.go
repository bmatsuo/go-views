package main

import (
	"flag"
)

type Options struct {
	Extension string
	Package   string
	Output    string
	Watch     bool
}

func parseOptions() (*Options, []string) {
	opt := new(Options)
	flag.StringVar(&opt.Extension, "ext", ".ghtml", "template file extension")
	flag.StringVar(&opt.Package, "pkg", "views", "view package name")
	flag.StringVar(&opt.Output, "o", "-", "output filename")
	flag.BoolVar(&opt.Watch, "w", false, "watch and recompile")
	flag.Parse()
	return opt, flag.Args()
}
