package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"

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
		atxCoinCmd,
		atxNonceCmd,
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
		flag.UnitFlag,
	},
	Action: atxLoad,
}

var atxCoinCmd = cli.Command{
	Name:  "coin",
	Usage: "load latest n ATXs per smesher.",
	Flags: []cli.Flag{
		flag.PathFlag,
	},
	Action: atxCoin,
}

var atxNonceCmd = cli.Command{
	Name:  "nonce",
	Usage: "load nonce of atx for a given epoch.",
	Flags: []cli.Flag{
		flag.PathFlag,
		flag.NodeFlag,
		flag.EpochFlag,
	},
	Action: atxNonce,
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

func sliceToCheckpointFmt(checkpoints []atxs.CheckpointAtx) (r []CheckpointFmt) {
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
	defer tx.Release()
	data, err := atxs.LatestN(tx, 10)
	if err != nil {
		return
	}

	var numUnit uint32
	for _, checkpoint := range data {
		numUnit += checkpoint.NumUnits
	}

	if ctx.Bool(flag.UnitFlag.Name) {
		fmt.Println("numUnit", numUnit)
		fmt.Println("power(T)", numUnit*64/1024)
		return
	}

	dataBytes, err := json.Marshal(sliceToCheckpointFmt(data))
	if err != nil {
		return
	}
	fmt.Println(string(dataBytes))

	return
}

type CoinbaseInfo struct {
	AddrCount int
	NumUnit   int
	Coinbase  string
}

func atxCoin(ctx *cli.Context) (err error) {
	sqlDB, err := sql.Open("file:" + ctx.String(flag.PathFlag.Name))
	if err != nil {
		return
	}
	defer sqlDB.Close()
	tx, err := sqlDB.Tx(context.Background())
	if err != nil {
		return
	}
	defer tx.Release()
	data, err := atxs.LatestN(tx, 10)
	if err != nil {
		return
	}

	fmt.Fprintf(os.Stderr, "#Smesher:\t%d\n", len(data))
	info := make(map[string]*CoinbaseInfo)
	smesherMap := make(map[string]int)

	for _, checkpoint := range data {
		if smesherMap[checkpoint.SmesherID.ShortString()] == 1 {
			// dataBytes, _ := json.Marshal(toCheckpointFmt(checkpoint))
			fmt.Fprintf(os.Stderr, "Dup Smesher:\t%s\n", checkpoint.SmesherID.String())
		}
		smesherMap[checkpoint.SmesherID.ShortString()] += 1
		cb, ok := info[checkpoint.Coinbase.String()]
		if !ok {
			cb = &CoinbaseInfo{}
			info[checkpoint.Coinbase.String()] = cb
		}
		cb.AddrCount += 1
		cb.NumUnit += int(checkpoint.NumUnits)
	}

	cbs := make([]*CoinbaseInfo, 0)
	for k, cb := range info {
		cb.Coinbase = k
		cbs = append(cbs, cb)
	}

	sort.Slice(cbs, func(i, j int) bool {
		return cbs[i].NumUnit > cbs[j].NumUnit
	})

	dataBytes, err := json.Marshal(cbs)
	if err != nil {
		return
	}
	fmt.Fprintf(os.Stderr, "#Coinbase:\t%d\n", len(cbs))
	fmt.Println(string(dataBytes))

	return
}

func atxNonce(ctx *cli.Context) (err error) {
	node := types.NodeID(types.HexToHash32(ctx.String(flag.NodeFlag.Name)))

	sqlDB, err := sql.Open("file:" + ctx.String(flag.PathFlag.Name))
	if err != nil {
		return
	}
	defer sqlDB.Close()

	tx, err := sqlDB.Tx(context.Background())
	if err != nil {
		return
	}
	defer tx.Release()

	nonce, err := atxs.VRFNonce(tx, node, types.EpochID(ctx.Int(flag.EpochFlag.Name)))
	if err != nil {
		return
	}

	fmt.Println("nonce", nonce)
	return
}
