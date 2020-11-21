package treedb

import (
	"errors"
	"fmt"
	"strings"
)

const (
	ExistingChild = 1 + iota
	ExistingSibling
	DoNothing
	NewChild
	NewSibling
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


func (parent *node) Add(pNode *node) (*node, int) {
	if parent.child == nil {
		parent.child = pNode
		return parent.child, NewChild
	} else {
		if parent.child.tag == pNode.tag {
			return parent.child, ExistingChild
		}
		var temp *node = parent
		for ; temp.sibling != nil; temp = temp.sibling {
			if temp.tag == pNode.tag {
				return temp, ExistingSibling
			}
		}
		temp.sibling = pNode
		return temp.sibling, NewSibling
	}
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
func (root *node) GenerateNodes( path string) {
	var parent *node = root
	words := strings.Split(path, "/")
	for _, word := range words {
		node := NewNode()
		node.SetTag(word)

		pos, retCode := parent.Add(node)
		switch retCode {
		case NewChild :
			parent = parent.child
		case NewSibling :
			parent = parent.sibling
		case ExistingChild :
			parent = parent.child
		case ExistingSibling: // 여러개의 sibling이면,
			parent = pos
		}
	}
}

func (root *node) Print() {
	root.print(0)
}

func (n *node) print(leftAlign int) {
	for i := 0; i < leftAlign; i++ {
		fmt.Printf("    ")
	}
	fmt.Println(n.tag)

	if n.child != nil {
		n.child.print(leftAlign+1)
	}
	if n.sibling != nil {
		n.sibling.print(leftAlign)
	}
}
