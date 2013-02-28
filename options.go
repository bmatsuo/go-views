package main

import (
	"flag"
)

type Options struct {
	Extension string
	Package   string
	Output    string
	Watch     bool
	HTML      bool
	templates string
	Debug     int
}

func parseOptions() (*Options, []string) {
	opt := new(Options)
	flag.StringVar(&opt.Extension, "ext", ".ghtml", "template file extension")
	flag.BoolVar(&opt.HTML, "html", true, "use html/template instead of text/template")
	flag.StringVar(&opt.Package, "pkg", "views", "view package name")
	flag.StringVar(&opt.Output, "o", "-", "output filename")
	flag.BoolVar(&opt.Watch, "w", false, "watch and recompile")
	flag.IntVar(&opt.Debug, "d", 0, "debug print level")
	flag.Parse()

	opt.templates = "text/template"
	if opt.HTML {
		opt.templates = "html/template"
	}
	return opt, flag.Args()
}
