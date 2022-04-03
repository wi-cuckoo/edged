package edged

import "github.com/urfave/cli"

var DebugFlag = cli.BoolFlag{
	EnvVar: "EDGED_DEBUG",
	Name:   "debug",
	Usage:  "enable running in debug mode",
}

var PubAddrFlag = cli.StringFlag{
	EnvVar: "EDGED_PUB_ADDR",
	Name:   "pub_addr",
	Value:  ":9527",
	Usage:  "tcp server address to publish message",
}

var SubAddrFlag = cli.StringFlag{
	EnvVar: "EDGED_SUB_ADDR",
	Name:   "sub_addr",
	Value:  ":4399",
	Usage:  "tcp server address to subscribe message",
}
