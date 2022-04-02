package utils

import (
	"fmt"
	"testing"
)

// 0 -> 1
// 0 -> 2
// 0 -> 3
// 1 -> 4
// 1 -> 5
// 4 -> 6
// 3 -> 7
// 6 -> 8
func TestTopoSort(t *testing.T) {
	root := NewNode("0")
	l1 := NewNode("1")
	l2 := NewNode("2")
	l3 := NewNode("3")
	l4 := NewNode("4")
	l5 := NewNode("5")
	l6 := NewNode("6")
	l7 := NewNode("7")
	l8 := NewNode("8")

	root.AppendChild(l1)
	root.AppendChild(l2)
	root.AppendChild(l3)
	l1.AppendChild(l4)
	l1.AppendChild(l5)
	l4.AppendChild(l6)
	l3.AppendChild(l7)
	l6.AppendChild(l8)

	fmt.Println(GetLevels(root))
}
