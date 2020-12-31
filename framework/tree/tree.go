package tree

// TODO: This is ane experimental package, the idea being you can add it in and
// add it to any struct making it a tree object with all the necssary tree
// functions to just work as a tree.

// TODO: Our tree library should have merkle functionality by default
// References:
//   https://github.com/armon/go-radix/blob/master/radix.go
//   https://github.com/xlab/treeprint/blob/master/treeprint.go

// https://github.com/cbergoon/merkletree

// https://github.com/emirpasic/gods/blob/master/examples/redblacktree/redblacktree.go

type Value interface{}

type MerkleHash string

type Tree interface {
	Add(v Value) Tree
	Remove(hash string) Tree
	Subtree(hash string) Tree

	Value() []Value

	Contains(v Value)

	//AddBranch

	Branch() Node

	Search(v Value) Node

	LastNode() Node

	IsNode(hash string) bool
	MerkleRoot() string
	Verify() bool

	Clear() Tree
}

type Node struct {
	Parent   *Node
	Children []*Node
}
