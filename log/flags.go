package log

import (
	"flag"
)

var (
	verbose bool
	debug   bool
	alerts  bool
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "Verbose enables detailed logging")
	flag.BoolVar(&debug, "debug", false, "Debug enables debug logging")
	flag.BoolVar(&alerts, "alerts", false, "Alerts enables the machine to speak important messages")
}
