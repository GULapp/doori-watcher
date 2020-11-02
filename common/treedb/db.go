package treedb

import "strings"

type node struct {
	tag 	string
	child   *node
	brother *node
	data    *interface{}
}

func NewNode() *node {
	return &node{"",nil, nil, nil}
}

func (n *node) AddChildNode(pNode *node) {
	n.child = pNode
}

func (n *node) AddBrotherNode(pNode *node) {
	n.brother = pNode
}

func (n *node) LinkDataTable(pData *interface{}) {
	n.data = pData
}

//재귀함수로 구현해서, node를 정리하도록 한다.
//todo : 만약에 이 treedb의 깊이가 너무 깊으면 재귀함수 스택이 너무 올라간다.
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
// /root/child/child/brother ... 형식의 문자열
func (root *node)GetDataTable(path string) (*interface{}, error) {
	words := strings.Split(path, "/")
	for i := range words {
		root.EqualTag(words[i])
	}
}

func (n *node)find(key string, depth int) bool {
	return true
}

func (n *node)EqualTag(value string) bool {
	return n.tag == value
}