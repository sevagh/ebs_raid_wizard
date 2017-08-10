package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/sevagh/goat/awsutil"
	"github.com/sevagh/goat/commands"
)

var VERSION string

func main() {
	usage := `goat - EC2/EBS utility

Usage:
  goat disk [--log-level=<log-level>] [--dry] [--debug]
  goat network [--log-level=<log-level>] [--dry] [--debug]
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

	log.Printf("%s", commands.DrawASCIIBanner("WELCOME TO GOAT", debug))
	log.Printf("%s", commands.DrawASCIIBanner("1: COLLECTING EC2 INFO", debug))
	ec2Instance := awsutil.GetEC2InstanceData()

	cmd := arguments["<command>"].(string)

	switch cmd {
	case "disk":
		commands.GoatDisk(ec2Instance, dryRun, debug)
	case "network":
		log.Fatalf("Network feature hasn't been implemented in Goat yet")
	}
}
