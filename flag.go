package edged

import "github.com/urfave/cli"

var DebugFlag = cli.BoolFlag{
	EnvVar: "EDGED_DEBUG",
	Name:   "debug",
	Usage:  "enable running in debug mode",
}

var TCPAddrFlag = cli.StringFlag{
	EnvVar: "EDGED_TCP_ADDR",
	Name:   "tcp_addr",
	Value:  ":9527",
	Usage:  "tcp server address to publish message",
}
