package treedb

import (
	"errors"
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

func (n *node) addChildNode(pNode *node) *node {
	if n.tag == pNode.tag {
		return n
	}
	var temp *node = n.sibling
	if n.child == nil {
		n.child = pNode
		temp = n.child
	} else {
		temp = n.sibling
		for temp.sibling != nil {
			temp=temp.sibling
		}
		temp.sibling = pNode
	}
	return temp
}

func (n *node) linkDataTable(pData *interface{}) {
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

	return findNode, nil
}

// path 는 root/node1/node2 형식값으로 된 문자열
func (root *node) GenerateNodes(path string) {

}

// 해당되는 Child node가 존재하면, (*node, true)
// 해당되는 Child node가 존재 안하면, parent node를 리턴 (*parent node, false)
func (n *node) GetChildNextNode(tag string) (*node, bool) {
	if n.tag == tag {
		return n, true
	}
	if n.child == nil{
		return n, false
	} else {
		return n.child.GetChildNextNode(tag)
	}
}

// 해당되는 Brother node가 존재하면, (*node, true)
// 해당되는 Brother node가 존재 안하면, last Brother node를 리턴(*last sibling node, false)
func (n *node) GetBrotherNextNode(tag string) (*node, bool) {
	if n.tag == tag {
		return n, true
	}
	if n.sibling == nil {
		return n, false
	} else {
		return n.sibling.GetBrotherNextNode(tag)
	}
}
