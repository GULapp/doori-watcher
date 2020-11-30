package treedb

import (
	"testing"
)

func TestDB(t *testing.T) {
	root := NewNode()
	root.GenerateNodes("root/a/b/c")
	root.GenerateNodes("root/a/c/d")
	root.GenerateNodes("root/a/d")
	root.GenerateNodes("root/aa/bb")
	root.GenerateNodes("root/aa/cc")

	if _, ret := root.Find("root/a/b/c");ret != nil {
		t.Fatalf("root/a/b/c")
	}
	if _, ret := root.Find("root/a/c/d");ret != nil {
		t.Fatalf("root/a/c/d")
	}
	if _, ret := root.Find("root/a/d");ret != nil {
		t.Fatalf("root/a/d")
	}
	if _, ret := root.Find("root/aa/bb");ret != nil {
		t.Fatalf("root/aa/bb")
	}
	if _, ret := root.Find("root/aa/cc");ret != nil {
		t.Fatalf("root/aa/cc")
	}
	if _, ret := root.Find("root/a/cc");ret == nil {
		t.Fatalf("root/aa/cc")
	}
	root.Print()
}

