package main

import (
	"fmt"
	"strings"
)

type TreeNode struct {
	ID       string
	Title    string
	Children []*TreeNode
}

func (t *TreeNode) GetByTitle(title string) *TreeNode {
	if t.Title == title {
		return t
	}
	for _, c := range t.Children {
		ret := c.GetByTitle(title)
		if ret != nil {
			return ret
		}
	}
	return nil
}

func (t *TreeNode) Print(order int) {
	if order == 0 {
		fmt.Println(t.Title)
	} else {
		var count int
		if order-1 > 0 {
			count = order - 1
		}
		fmt.Printf("%s%s%s\n", strings.Repeat("    ", count), "|---", t.Title)
		// fmt.Printf("%s%s%s: %s\n", strings.Repeat("    ", count), "|---", t.Title, t.ID)
	}
	for _, c := range t.Children {
		c.Print(order + 1)
	}
}
