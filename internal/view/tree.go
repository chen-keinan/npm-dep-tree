package view

import (
	"fmt"
	"github.com/chen-keinan/npm-dep-tree/pkg/model"
	"github.com/tufin/asciitree"
)

//BuildTreeView build dependencies tree view
func BuildTreeView(tv *model.DependencyTree, tree *asciitree.Tree, path string) {
	current := fmt.Sprintf("%s/%s:%s", path, tv.Name, tv.Version)
	tree.Add(current)
	for _, node := range tv.Dependencies {
		BuildTreeView(node, tree, current)
	}
}
