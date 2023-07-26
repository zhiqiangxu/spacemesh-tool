package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli"
	"github.com/zhiqiangxu/spacemesh-tool/cmd"
)

func setupAPP() *cli.App {
	app := cli.NewApp()
	app.Usage = "Spacemesh Cli"
	app.Copyright = "Copyright in 2020"
	app.Commands = []cli.Command{
		cmd.AtxCmd,
		cmd.NiPostCmd,
		cmd.MerkleCmd,
	}
	app.Flags = []cli.Flag{}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func main() {
	if err := setupAPP().Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
