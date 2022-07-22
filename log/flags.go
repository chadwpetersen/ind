package log

import (
	"flag"
)

var (
	verbose bool
	debug   bool
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "Verbose enables detailed logging")
	flag.BoolVar(&debug, "debug", false, "Debug enables debug logging")
}
