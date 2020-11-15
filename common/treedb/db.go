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

// -1, add 못함
// 0, child add
// 1, sibling add
func (n *node) AddNode(pNode *node) int {
	if pNode == nil {
		return -1
	}

	if n.tag == pNode.tag {
		return -1
	}

	if n.child == nil {
		n.child = pNode
		return 0
	} else {
		temp := n.sibling
		for temp != nil {
			temp = temp.sibling
		}
		temp = pNode
	}
	return 1
}

func (n *node) LinkDataTable(pData *interface{}) {
	n.data = pData
}

//재귀함수로 구현해서, node를 정리하도록 한다.
//todo : 만약에 이 treedb의 깊이가 너무 깊으면 재귀함수 스택수가 증가한다.
func (n *node) DestoryNode() {
	if n.child != nil {
		n.child.DestoryNode()
	}
	if n.sibling != nil {
		n.sibling.DestoryNode()
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
	var iter *node = root
	words := strings.Split(path, "/")
	for _, word := range words {
		node := NewNode()
		node.SetTag(word)

		if ret := iter.AddNode(node); ret == 0 {
			iter = iter.child
		} else if ret == 1{

		} else {
			node.DestoryNode()
		}
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

	n.child.print(leftAlign+1)
	n.sibling.print(leftAlign)
}
