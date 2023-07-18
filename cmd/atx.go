package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/zhiqiangxu/spacemesh-tool/cmd/flag"
)

// AtxCmd ...
var AtxCmd = cli.Command{
	Name:  "atx",
	Usage: "atx actions",
	Subcommands: []cli.Command{
		atxMetaCmd,
	},
}

var atxMetaCmd = cli.Command{
	Name:  "meta",
	Usage: "decode atxid from postdata_metadata.json into hex",
	Flags: []cli.Flag{
		flag.PathFlag},
	Action: atxMeta,
}

type PostMetadata struct {
	NodeId          []byte
	CommitmentAtxId []byte

	LabelsPerUnit uint64
	NumUnits      uint32
	MaxFileSize   uint64
	Nonce         *uint64 `json:",omitempty"`
	LastPosition  *uint64 `json:",omitempty"`
}

func atxMeta(ctx *cli.Context) (err error) {

	metaData, err := os.ReadFile(ctx.String(flag.PathFlag.Name))
	if err != nil {
		return
	}

	var v PostMetadata
	err = json.Unmarshal(metaData, &v)
	if err != nil {
		return
	}

	fmt.Println("CommitmentAtxId", hex.EncodeToString(v.CommitmentAtxId))
	fmt.Println("NodeId", hex.EncodeToString(v.NodeId))

	return
}
