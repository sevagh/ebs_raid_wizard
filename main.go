package main

import (
	"bufio"
	"fmt"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"

	"github.com/sevagh/goat/awsutil"
)

var VERSION string

func main() {
	usage := `goat - EC2/EBS utility

Usage:
  goat [--log-level=<log-level>] [--dry] [--debug]
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

	log.Printf("%s", drawASCIIBanner("WELCOME TO GOAT", debug))

	log.Printf("%s", drawASCIIBanner("1: COLLECTING EC2 INFO", debug))
	ec2Instance := awsutil.GetEC2InstanceData()

	log.Printf("%s", drawASCIIBanner("2: COLLECTING EBS INFO", debug))
	ebsVolumes := awsutil.MapEbsVolumes(&ec2Instance)

	log.Printf("%s", drawASCIIBanner("3: ATTACHING EBS VOLS", debug))
	ebsVolumes = awsutil.AttachEbsVolumes(ec2Instance, ebsVolumes, dryRun)

	log.Printf("%s", drawASCIIBanner("4: MOUNTING ATTACHED VOLS", debug))

	if len(ebsVolumes) == 0 {
		log.Warn("Empty vols, nothing to do")
		os.Exit(0)
	}

	for volName, vols := range ebsVolumes {
		PrepAndMountDrives(volName, vols, dryRun)
	}
}

func drawASCIIBanner(headLine string, debug bool) string {
	if debug {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Press enter to continue: ")
		reader.ReadString('\n')
	}

	return fmt.Sprintf("\n%[1]s\n# %[2]s #\n%[1]s\n",
		strings.Repeat("#", len(headLine)+4),
		headLine)
}
