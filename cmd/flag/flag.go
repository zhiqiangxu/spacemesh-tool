package flag

import "github.com/urfave/cli"

var PathFlag = cli.StringFlag{
	Name:  "path",
	Usage: "specify file path",
}
