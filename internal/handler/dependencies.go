package handler

import (
	"fmt"
	"github.com/chen-keinan/npm-dep-tree/internal/service"
	"github.com/chen-keinan/npm-dep-tree/internal/view"
	"github.com/chen-keinan/npm-dep-tree/pkg/model"
	"github.com/gorilla/mux"
	"github.com/tufin/asciitree"
	"go.uber.org/zap"
	"net/http"
)

//Dependencies handler object
type Dependencies struct {
	log        *zap.Logger
	depService service.Dep
}

//NewDependenciesHandler return new dependencies handler instance
func NewDependenciesHandler(zlog *zap.Logger, depService service.Dep) *Dependencies {
	return &Dependencies{log: zlog, depService: depService}
}

//ResolveDependencies handler by package name and version
func (handler Dependencies) ResolveDependencies(w http.ResponseWriter, req *http.Request) {
	name := mux.Vars(req)["name"]
	version := mux.Vars(req)["version"]
	if len(name) == 0 || len(version) == 0 {
		handler.log.Error("name or version are missing from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rootTree := &model.DependencyTree{Name: name, Version: version, Dependencies: []*model.DependencyTree{}}
	err := handler.depService.ResolveDependencies(rootTree)
	if err != nil {
		handler.log.Error(fmt.Sprintf("failed to fetch dependencies tree for pkg:%s and version:%s", name, version))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	tree := asciitree.Tree{}
	// build tree view
	view.BuildTreeView(rootTree, &tree, "")
	// render tree to ui
	tree.Fprint(w, false, "")
}
