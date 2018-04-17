package main // import "github.com/sevagh/goat"

import (
	"fmt"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"os"
)

//Goat version substituted by the Makefile
var VERSION string

func main() {
	usage := `goat - EC2/EBS attach utility

Usage:
  goat ebs [--log-level=<log-level>] [--dry] [--debug]
  goat eni [--log-level=<log-level>] [--dry] [--debug]
  goat -h | --help
  goat --version

Options:
  --log-level=<level>  Log level (debug, info, warn, error, fatal) [default: info]
  --dry                Dry run
  --debug              Interactive prompts to continue between phases
  -h --help            Show this screen.
  --version            Show version.`
	arguments, _ := docopt.Parse(usage, nil, true, fmt.Sprintf("goat %s", VERSION), false)

	log.SetOutput(os.Stderr)
	logLevel := arguments["--log-level"].(string)
	if level, err := log.ParseLevel(logLevel); err != nil {
		log.Fatalf("%v", err)
	} else {
		log.SetLevel(level)
	}

	log.SetFormatter(&log.TextFormatter{})

	dryRun := arguments["--dry"].(bool)
	debug := arguments["--debug"].(bool)

	if arguments["ebs"].(bool) {
		log.Printf("Running goat for EBS")
		GoatEbs(dryRun, debug)
	} else if arguments["eni"].(bool) {
		log.Printf("Running goat for ENI")
		GoatEni(dryRun, debug)
	}
}
