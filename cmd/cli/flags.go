package main

import "flag"

type Flags struct {
	forceResetTemplates bool
}

var flags Flags

func init() {
	flag.BoolVar(&flags.forceResetTemplates, "forceDefault", false, "Reset default templates")
	flag.BoolVar(&flags.forceResetTemplates, "f", false, "Reset default templates")
	flag.Parse()
}
