package treedb

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	ExistingChild = 1 + iota
	ExistingSibling
	NewChild
	NewSibling
)

type node struct {
	tag     string
	child   *node
	sibling *node
	data    interface{}
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
		var temp **node = &parent.sibling
		for ; *temp != nil; temp = &(*temp).sibling {
			if (*temp).tag == pNode.tag {
				return *temp, ExistingSibling
			}
		}
		*temp = pNode
		return *temp, NewSibling
	}
}

func (parent *node) find(tag string) (*node, bool) {
	if parent.child != nil {
		if parent.child.tag == tag {
			return parent.child, true
		}
	}
	var temp **node = &parent.sibling
	for ; *temp != nil; temp = &(*temp).sibling {
		if (*temp).tag == tag {
			return *temp, true
		}
	}
	return nil, false
}

// 해당 노드에 노드가 사라지기전까지의 저장할 데이타구조를 정의한다.
// 저장할 데이터가 어느형식의 포인터가 되어야 한다. heap 메모리에 저장되어야 함으로
// 내부적으로 포인터형인지 체크함
func (n *node) LinkDataTable(pData interface{}) bool {
	if reflect.ValueOf(pData).Kind() != reflect.Ptr {
		return false
	}
	n.data = pData
	return true
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

func (root *node) Find(path string) (*node, error) {
	var parent *node = root
	var isFound bool = false
	words := strings.Split(path, "/")
	for _, tag := range words { if parent, isFound = parent.find(tag); isFound == true {
			continue
		} else {
			return nil, errors.New("not found")
		}
	}
	return parent, nil
}

func (root *node) FindFromArgs(words ...string) (*node, error) {
	var parent *node = root
	var isFound bool = false
	for _, tag := range words {
		if parent, isFound = parent.find(tag); isFound == true {
			continue
		} else {
			return nil, errors.New("not found")
		}
	}
	return parent, nil
}



// path 는 root/node1/node2 형식값으로 된 문자열
func (root *node) GenerateNodes(path string) *node {
	var parent *node = root
	words := strings.Split(path, "/")
	for _, word := range words {
		node := NewNode()
		node.SetTag(word)

		pos, retCode := parent.Add(node)
		switch retCode {
		case NewChild:
			parent = parent.child
		case NewSibling:
			parent = parent.sibling
		case ExistingChild:
			parent = parent.child
		case ExistingSibling: // 여러개의 sibling이면,
			parent = pos
		}
	}
	return parent
}

func (root *node) GenerateNodesFromArgs(words ...string) *node {
	var parent *node = root
	for _, word := range words {
		node := NewNode()
		node.SetTag(word)

		pos, retCode := parent.Add(node)
		switch retCode {
		case NewChild:
			parent = parent.child
		case NewSibling:
			parent = parent.sibling
		case ExistingChild:
			parent = parent.child
		case ExistingSibling: // 여러개의 sibling이면,
			parent = pos
		}
	}
	return parent
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
		n.child.print(leftAlign + 1)
	}
	if n.sibling != nil {
		n.sibling.print(leftAlign + 1)
	}
}

func (n node) Tag() string {
	return n.tag
}

func (n *node) GetDataTable() interface{} {
	return n.data
}