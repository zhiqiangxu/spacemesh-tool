package flag

import "github.com/urfave/cli"

var PathFlag = cli.StringFlag{
	Name:  "path",
	Usage: "specify file path",
}

var UnitFlag = cli.BoolFlag{
	Name:  "unit",
	Usage: "specify num unit",
}

var NodeFlag = cli.StringFlag{
	Name:  "node",
	Usage: "specify node",
}

var EpochFlag = cli.IntFlag{
	Name:  "epoch",
	Usage: "specify epoch",
}
