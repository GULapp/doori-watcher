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

// -2, failed to add
// -1, existing Node
// 0, add
func (n *node) Add(pNode *node) int {
	if pNode == nil || pNode.tag == "" {
		return -2
	}
	if n.tag == pNode.tag {
		return -1
	}
	if n.addChild(pNode) != nil {
		return -1
	} else {
		if n.addSibling(pNode) != nil {
			return -1
		}
	}
	return 0
}

func (n *node) addChild(pNode *node) error {
	if n.child == nil {
		n.child = pNode
		return nil
	}
	if n.child.tag == pNode.tag {
		return n.addSibling(pNode)
	}
	return nil
}

func (n *node) addSibling(pNode *node) error {
	temp := &n.sibling
	for *temp != nil {
		temp = &(*temp).sibling
		if (*temp).tag == pNode.tag {
			return errors.New("sibling is duplicated")
		}
		*temp = pNode
	}
	return nil
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
	var iter *node = root
	words := strings.Split(path, "/")
	for _, word := range words {
		node := NewNode()
		node.SetTag(word)

		if ret := iter.Add(node); ret == 0 {
			iter = iter.child
		} else {
			node.DestroysNode()
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

	if n.tag == "" {
		n.child.print(leftAlign)
	} else {
		n.child.print(leftAlign+1)
	}
	n.sibling.print(leftAlign)
}
