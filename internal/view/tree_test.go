package view

import (
	"bytes"
	"github.com/chen-keinan/npm-dep-tree/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/tufin/asciitree"
	"testing"
)

const (
	treeView = `└ 
  └ a:1.0.1
    ├ b:2.0.0
    │ └ c:3.0.0
    └ k:6.0.0
`
	treeView2 = `└ 
  └ a:1.0.1
    ├ k:6.0.0
    └ b:2.0.0
      └ c:3.0.0
`
	singleView = `└ 
  └ a:1.0.1
`
)

func Test_BuildView(t *testing.T) {
	p := &model.DependencyTree{Name: "a", Version: "1.0.1",
		Dependencies: []*model.DependencyTree{{Name: "b", Version: "2.0.0",
			Dependencies: []*model.DependencyTree{{Name: "c", Version: "3.0.0",
				Dependencies: []*model.DependencyTree{}}}}, {Name: "k", Version: "6.0.0"}}}
	tree := asciitree.Tree{}
	BuildTreeView(p, &tree, "")
	buf := new(bytes.Buffer)
	tree.Fprint(buf, false, "")
	s := buf.String()
	assert.True(t, s == treeView || s == treeView2)
}

func Test_BuildSinglePkgView(t *testing.T) {
	p := &model.DependencyTree{Name: "a", Version: "1.0.1"}
	tree := asciitree.Tree{}
	BuildTreeView(p, &tree, "")
	buf := new(bytes.Buffer)
	tree.Fprint(buf, false, "")
	s := buf.String()
	assert.Equal(t, s, singleView)
}
