package flag

import "github.com/urfave/cli"

var PathFlag = cli.StringFlag{
	Name:  "path",
	Usage: "specify file path",
}

var NumUnit = cli.BoolFlag{
	Name:  "unit",
	Usage: "specify num unit",
}
