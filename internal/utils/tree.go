package utils

import "fmt"

type Node struct {
	Value    interface{}
	Parent   *Node
	Children []*Node
}

func NewNode(v interface{}) *Node {
	return &Node{Value: v}
}

func IsRoot(n *Node) bool {
	return n.Parent == nil
}

func (parent *Node) AppendChild(n *Node) *Node {
	if parent == nil || n == nil || !IsRoot(n) {
		return nil
	}

	n.Parent = parent
	parent.Children = append(parent.Children, n)

	return n
}

func GetLevels(n *Node) [][]string {
	levels := [][]string{}
	level := n.Children

	for len(level) > 0 {
		fmt.Println(GetStringValues(level))
		levels = append(levels, GetStringValues(level))

		children := []*Node{}
		for _, leaf := range level {
			children = append(children, leaf.Children...)
		}

		level = children
	}

	return levels
}

func GetStringValues(nodes []*Node) []string {
	vals := []string{}
	for _, n := range nodes {
		vals = append(vals, n.Value.(string))
	}

	return vals
}
