package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

// MerkleCmd ...
var MerkleCmd = cli.Command{
	Name:  "merkle",
	Usage: "merkle actions",
	Subcommands: []cli.Command{
		merkleIsAncestorCmd,
		merkleSiblingCmd,
		merkleParentCmd,
		merkleIsRightCmd,
		merkleChildrenCmd,
	},
}

var merkleIsAncestorCmd = cli.Command{
	Name:   "is_ancestor",
	Usage:  "check whether one node is ancestor of another, format is_ancestor ancestor descendant, node format: height:index",
	Flags:  []cli.Flag{},
	Action: merkleIsAncestor,
}

var merkleSiblingCmd = cli.Command{
	Name:   "sibling",
	Usage:  "print the sibling node",
	Flags:  []cli.Flag{},
	Action: merkleSibling,
}

var merkleParentCmd = cli.Command{
	Name:   "parent",
	Usage:  "print the parent node",
	Flags:  []cli.Flag{},
	Action: merkleParent,
}

var merkleIsRightCmd = cli.Command{
	Name:   "is_right",
	Usage:  "check whether is right node",
	Flags:  []cli.Flag{},
	Action: merkleIsRight,
}

var merkleChildrenCmd = cli.Command{
	Name:   "children",
	Usage:  "print node's left and right child",
	Flags:  []cli.Flag{},
	Action: merkleChildren,
}

type Position struct {
	Index  uint64
	Height uint64
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Height, p.Index)
}

func (p Position) sibling() Position {
	return Position{
		Index:  p.Index ^ 1,
		Height: p.Height,
	}
}

func (p Position) isAncestorOf(other Position) bool {
	if p.Height < other.Height {
		return false
	}
	return p.Index == (other.Index >> (p.Height - other.Height))
}

func (p Position) isRightSibling() bool {
	return p.Index%2 == 1
}

func (p Position) parent() Position {
	return Position{
		Index:  p.Index >> 1,
		Height: p.Height + 1,
	}
}

func (p Position) leftChild() Position {
	return Position{
		Index:  p.Index << 1,
		Height: p.Height - 1,
	}
}

func (p Position) rightChild() Position {
	return p.leftChild().sibling()
}

func parseNode(nodeStr string) (p Position, err error) {
	parts := strings.Split(nodeStr, ":")
	if len(parts) != 2 {
		err = fmt.Errorf("invalid nodeStr:%v", nodeStr)
		return
	}

	p.Height, err = strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return
	}
	p.Index, err = strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return
	}
	return
}

func merkleChildren(ctx *cli.Context) (err error) {
	if len(ctx.Args()) != 1 {
		err = fmt.Errorf("invalid input")
		return
	}

	node, err := parseNode(ctx.Args()[0])
	if err != nil {
		return
	}

	fmt.Println(node.leftChild().String(), node.rightChild().String())
	return
}

func merkleIsRight(ctx *cli.Context) (err error) {
	if len(ctx.Args()) != 1 {
		err = fmt.Errorf("invalid input")
		return
	}

	node, err := parseNode(ctx.Args()[0])
	if err != nil {
		return
	}

	fmt.Println(node.isRightSibling())
	return
}

func merkleParent(ctx *cli.Context) (err error) {
	if len(ctx.Args()) != 1 {
		err = fmt.Errorf("invalid input")
		return
	}

	node, err := parseNode(ctx.Args()[0])
	if err != nil {
		return
	}

	fmt.Println(node.parent().String())
	return
}

func merkleSibling(ctx *cli.Context) (err error) {
	if len(ctx.Args()) != 1 {
		err = fmt.Errorf("invalid input")
		return
	}

	node, err := parseNode(ctx.Args()[0])
	if err != nil {
		return
	}

	fmt.Println(node.sibling().String())
	return
}

func merkleIsAncestor(ctx *cli.Context) (err error) {
	if len(ctx.Args()) != 2 {
		err = fmt.Errorf("invalid input")
		return
	}

	ancestor, err := parseNode(ctx.Args()[0])
	if err != nil {
		return
	}
	descendant, err := parseNode(ctx.Args()[1])
	if err != nil {
		return
	}

	fmt.Println(ancestor.isAncestorOf(descendant))
	return
}
