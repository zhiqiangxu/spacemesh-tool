package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/sql"
	"github.com/spacemeshos/go-spacemesh/sql/atxs"
	"github.com/urfave/cli"
	"github.com/zhiqiangxu/spacemesh-tool/cmd/flag"
)

// AtxCmd ...
var AtxCmd = cli.Command{
	Name:  "atx",
	Usage: "atx actions",
	Subcommands: []cli.Command{
		atxMetaCmd,
		atxLoadCmd,
	},
}

var atxMetaCmd = cli.Command{
	Name:  "meta",
	Usage: "decode atxid from postdata_metadata.json into hex",
	Flags: []cli.Flag{
		flag.PathFlag},
	Action: atxMeta,
}

var atxLoadCmd = cli.Command{
	Name:  "load",
	Usage: "load latest n ATXs per smesher.",
	Flags: []cli.Flag{
		flag.PathFlag,
		flag.NumUnit,
	},
	Action: atxLoad,
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

type CheckpointFmt struct {
	ID             string
	Epoch          types.EpochID
	CommitmentATX  string
	VRFNonce       types.VRFPostIndex
	NumUnits       uint32
	BaseTickHeight uint64
	TickCount      uint64
	SmesherID      string
	Sequence       uint64
	Coinbase       string
}

func toCheckpointFmt(checkpoints []atxs.CheckpointAtx) (r []CheckpointFmt) {
	for _, checkpoint := range checkpoints {
		r = append(r, CheckpointFmt{
			ID:             checkpoint.ID.ShortString(),
			Epoch:          checkpoint.Epoch,
			CommitmentATX:  checkpoint.CommitmentATX.ShortString(),
			VRFNonce:       checkpoint.VRFNonce,
			NumUnits:       checkpoint.NumUnits,
			BaseTickHeight: checkpoint.BaseTickHeight,
			TickCount:      checkpoint.TickCount,
			SmesherID:      checkpoint.SmesherID.ShortString(),
			Sequence:       checkpoint.Sequence,
			Coinbase:       checkpoint.Coinbase.String(),
		})
	}
	return
}

func atxLoad(ctx *cli.Context) (err error) {
	sqlDB, err := sql.Open("file:" + ctx.String(flag.PathFlag.Name))
	if err != nil {
		return
	}
	defer sqlDB.Close()
	tx, err := sqlDB.Tx(context.Background())
	if err != nil {
		return
	}
	data, err := atxs.LatestN(tx, 10)
	if err != nil {
		return
	}

	var numUnit uint32
	for _, checkpoint := range data {
		numUnit += checkpoint.NumUnits
	}

	if ctx.Bool(flag.NumUnit.Name) {
		fmt.Println("numUnit", numUnit)
		fmt.Println("power(T)", numUnit*64/1024)
		return
	}

	dataBytes, err := json.Marshal(toCheckpointFmt(data))
	if err != nil {
		return
	}
	fmt.Println(string(dataBytes))

	return
}
