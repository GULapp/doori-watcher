package treedb

import (
	"errors"
	"fmt"
	"strings"
)

type node struct {
	tag     string
	child   *node
	sibling *node
	data    *interface{}
}

func NewNode() *node {
	return &node{"", nil, nil, nil}
}

func (n *node) SetTag(tag string) error {
	if tag == "" {
		return errors.New("tag is nil")
	}
	n.tag = tag
	return nil
}

// -1, existing Node
// 0, add
// 이 함수의 호출하는 대상은 parent 형의 노드여야 한다.
func (parent *node) Add(pNode *node) int {
	if parent.child == nil {
		parent.child = pNode
	} else {
		if parent.child.tag == pNode.tag {
			return -1
		}
		var temp *node = parent.child
		for ; temp.sibling != nil; temp = temp.sibling {
			if temp.tag == pNode.tag {
				return -1
			}
		}
		temp.sibling = pNode
	}
	return 0
}

func (n *node) LinkDataTable(pData *interface{}) {
	n.data = pData
}

//재귀함수로 구현해서, node를 정리하도록 한다.
//todo : 만약에 이 treedb의 깊이가 너무 깊으면 재귀함수 스택수가 증가한다.
func (n *node) DestroysNode() {
	if n.child != nil {
		n.child.DestroysNode()
	}
	if n.sibling != nil {
		n.sibling.DestroysNode()
	}
	n.child = nil
	n.sibling = nil
	n.data = nil
}

func (pN *node) Find(tag string) (*node, error) {
	if pN == nil {
		return nil, errors.New("node is nil")
	}
	if pN.tag == tag {
		return pN, nil
	}
	if findNode, err := pN.child.Find(tag); err != nil {
		return findNode, nil
	}

	if findNode, err := pN.sibling.Find(tag); err != nil {
		return findNode, nil
	}
	return nil, errors.New("not found")
}

// path 는 root/node1/node2 형식값으로 된 문자열
func (root *node) GenerateNodes(path string) {
	var parent *node = root
	words := strings.Split(path, "/")
	for _, word := range words {
		node := NewNode()
		node.SetTag(word)

		parent.Add(node)
		parent = parent.child
	}
}

func (root *node) Print() {
	root.print(0)
}

func (n *node) print(leftAlign int) {
	if n == nil {
		return
	}
	for i := 0; i < leftAlign; i++ {
		fmt.Printf("    ")
	}
	fmt.Print(n.tag)
	fmt.Println("")

	if n.tag == "" {
		n.child.print(leftAlign)
	} else {
		n.child.print(leftAlign + 1)
	}
	n.sibling.print(leftAlign)
}
