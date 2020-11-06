package treedb

import (
	"fmt"
	"strings"
)

type node struct {
	tag     string
	child   *node
	brother *node
	data    *interface{}
}

func NewNode() *node {
	return &node{"", nil, nil, nil}
}

func (n *node) addChildNode(pNode *node) {
	n.child = pNode
}

func (n *node) addBrotherNode(pNode *node) {
	n.brother = pNode
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
	if n.brother != nil {
		n.brother.DestoryNode()
	}
	n.child = nil
	n.brother = nil
	n.data = nil
}

// /site/domain/server/cpu/... 형식의 문자열
func (root *node) GetDataTable(path string) *interface{} {
	words := strings.Split(path, "/")
	var pNode = root
	for i := range words {
		tag := words[i]

		if pNode, status := pNode.GetChildNextNode(tag); status == false {
			pNode.child= NewNode()
			pNode.child.tag = tag
			pNode = pNode.child
			fmt.Println("New V --->:", tag)
		}

		if pNode, status := pNode.GetBrotherNextNode(tag); status == false {
			pNode.brother = NewNode()
			pNode.brother.tag = tag
			fmt.Println("New --->:", tag)
		}
	}
	return pNode.data
}

// 해당되는 Brother node가 존재하면, (*node, true)
// 해당되는 Brother node가 존재 안하면, last Brother node를 리턴(*last brother node, false)
func (n *node) GetBrotherNextNode(tag string) (*node, bool) {
	if n.tag == tag {
		return n, true
	}
	if n.brother == nil {
		fmt.Println("--->/:", tag)
		return n, false
	} else {
		return n.brother.GetBrotherNextNode(tag)
	}
}

// 해당되는 Child node가 존재하면, (*node, true)
// 해당되는 Child node가 존재 안하면, parent node를 리턴 (*parent node, false)
func (n *node) GetChildNextNode(tag string) (*node, bool) {
	if n.tag == tag {
		return n, true
	}
	if n.child == nil{
		fmt.Println("V /:", tag)
		return n, false
	} else {
		return n.child.GetChildNextNode(tag)
	}
}
