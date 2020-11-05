package treedb

import (
	"strings"
)

type node struct {
	depth	uint
	tag 	string
	child   *node
	brother *node
	data    *interface{}
}

func NewNode() *node {
	return &node{0,"",nil, nil, nil}
}

func (n *node) addChildNode(pNode *node) {
	n.child = pNode
	n.child.depth = n.depth+1
}

func (n *node) addBrotherNode(pNode *node) {
	n.brother = pNode
	n.brother.depth = n.depth
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

// /depth1/depth2/depth3/depth4
// /site/domain/server/cpu/... 형식의 문자열
func (root *node)GetDataTable(path string) (*interface{}, error) {
	words := strings.Split(path, "/")
	for i := range words {
		if root.EqualTag(words[i]) {

		} else {

		}
	}
	return nil, nil
}

func (n *node)NextNode(depth uint,tag string) (*node, error) {
	return nil, nil
}

func (n *node)EqualTag(value string) bool {
	return n.tag == value
}
