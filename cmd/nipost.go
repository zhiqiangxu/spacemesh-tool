package cmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc64"
	"io"
	"os"

	"github.com/spacemeshos/go-scale"
	"github.com/spacemeshos/go-spacemesh/codec"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/urfave/cli"
	"github.com/zhiqiangxu/spacemesh-tool/cmd/flag"
)

// NiPostCmd ...
var NiPostCmd = cli.Command{
	Name:  "nipost",
	Usage: "nipost actions",
	Subcommands: []cli.Command{
		niPostChallengeCmd,
		niPostBuilderStateCmd,
	},
}

var niPostChallengeCmd = cli.Command{
	Name:  "challenge",
	Usage: "decode NIPostChallenge from nipost_challenge.bin",
	Flags: []cli.Flag{
		flag.PathFlag},
	Action: niPostChallenge,
}

var niPostBuilderStateCmd = cli.Command{
	Name:  "bs",
	Usage: "decode NIPostBuilderState from nipost_builder_state.bin",
	Flags: []cli.Flag{
		flag.PathFlag},
	Action: niPostBuilderState,
}

func load(filename string, dst scale.Decodable) error {
	data, err := read(filename)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	if err := codec.Decode(data, dst); err != nil {
		return fmt.Errorf("decoding: %w", err)
	}
	return nil
}

func read(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %w", path, err)
	}

	defer file.Close()

	fInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info %s: %w", path, err)
	}
	if fInfo.Size() < crc64.Size {
		return nil, fmt.Errorf("file %s is too small", path)
	}

	data := make([]byte, fInfo.Size()-crc64.Size)
	checksum := crc64.New(crc64.MakeTable(crc64.ISO))
	if _, err = io.TeeReader(file, checksum).Read(data); err != nil {
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}

	saved := make([]byte, crc64.Size)
	if _, err = file.Read(saved); err != nil {
		return nil, fmt.Errorf("read checksum %s: %w", path, err)
	}

	savedChecksum := binary.BigEndian.Uint64(saved)

	if savedChecksum != checksum.Sum64() {
		return nil, fmt.Errorf(
			"wrong checksum 0x%X, computed 0x%X", savedChecksum, checksum.Sum64())
	}

	return data, nil
}

func niPostChallenge(ctx *cli.Context) (err error) {

	var ch types.NIPostChallenge
	if err := load(ctx.String(flag.PathFlag.Name), &ch); err != nil {
		return fmt.Errorf("decoding: %w", err)
	}

	resultBytes, _ := json.Marshal(ch)
	fmt.Println("NIPostChallenge", string(resultBytes))
	fmt.Println()
	if ch.CommitmentATX != nil {
		fmt.Println("CommitmentATX", ch.CommitmentATX.Hash32().Hex())
	}

	fmt.Println("PositioningATX", ch.PositioningATX.Hash32().Hex())
	fmt.Println("PrevATXID", ch.PrevATXID.Hash32().Hex())
	return
}

func niPostBuilderState(ctx *cli.Context) (err error) {

	var ch types.NIPostBuilderState
	if err := load(ctx.String(flag.PathFlag.Name), &ch); err != nil {
		return fmt.Errorf("decoding: %w", err)
	}

	resultBytes, _ := json.Marshal(ch)
	fmt.Println("NIPostBuilderState", string(resultBytes))

	return
}
