package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/wi-cuckoo/edged"
)

var (
	// Revision ...
	revision string
	// Version ...
	version string
)

func main() {
	app := cli.NewApp()
	app.Name = "Edged"
	app.Version = version
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version=%s revision=%s\n", c.App.Version, revision)
	}
	app.Usage = "IOT broker with MQTT protocol"
	app.Action = setup
	app.Flags = []cli.Flag{
		edged.DebugFlag,
		edged.PubAddrFlag,
		edged.SubAddrFlag,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	fmt.Fprint(os.Stdout, edged.Banner)
}
