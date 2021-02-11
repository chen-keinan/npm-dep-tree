package service

import (
	"encoding/json"
	"fmt"
	"github.com/chen-keinan/npm-dep-tree/internal/cache"
	"github.com/chen-keinan/npm-dep-tree/pkg/model"
	"github.com/rcrowley/go-metrics"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

//NpmRegistry npm registry url
const NpmRegistry = "https://registry.npmjs.org/"

//Dep is interface for dependencies resolver service
//dependencies.go
//go:generate mockgen -destination=../mocks/mock_Dep.go -package=mocks . Dep
type Dep interface {
	getNextDependency(pkgName string, pkgVersion string) (*model.NpmDependency, error)
	fetchFromRegistry(pkgName string, pkgVersion string) (*model.NpmDependency, error)
	ResolveDependencies(rootTree *model.DependencyTree) error
	getLru() *cache.Lru
}

//Dependencies service struct
type Dependencies struct {
	log *zap.Logger
	lru *cache.Lru
}

//NewDependencies create dependencies service instance
func NewDependencies(zlog *zap.Logger, l *cache.Lru) Dep {
	return &Dependencies{log: zlog, lru: l}
}

//ResolveDependencies resolve npm package dependency by name and version
func (d Dependencies) ResolveDependencies(rootTree *model.DependencyTree) error {
	d.log.Debug(fmt.Sprintf("fetching dependencies for pkg name %s and version %s", rootTree.Name, rootTree.Version))
	npmDep, err := d.getNextDependency(rootTree.Name, rootTree.Version)
	if err != nil {
		return err
	}
	for name, version := range npmDep.Dependencies {
		dep := &model.DependencyTree{Name: name, Version: strings.Trim(version, "^"), Dependencies: []*model.DependencyTree{}}
		err := d.ResolveDependencies(dep)
		if err == nil {
			rootTree.Dependencies = append(rootTree.Dependencies, dep)
		}
	}
	return nil
}

func (d Dependencies) getNextDependency(pkgName string, pkgVersion string) (*model.NpmDependency, error) {
	// check if npm pkg exist in cache
	val, ok := d.lru.Get(fmt.Sprintf("%s:%s", pkgName, pkgVersion))
	if ok {
		d.log.Debug(fmt.Sprintf("resolving dependencies from cache for pkg name %s and version %s", pkgName, pkgVersion))
		return val.(*model.NpmDependency), nil
	}
	t := metrics.GetOrRegisterTimer("fetch.registry.dependency.latency", nil)
	var npmDep *model.NpmDependency
	var err error
	t.Time(func() {
		// if not fetch it from npm registry
		npmDep, err = d.fetchFromRegistry(pkgName, pkgVersion)
	})
	t.Update(47)
	if err != nil {
		return npmDep, err
	}
	d.log.Debug(fmt.Sprintf("resolving dependencies from registry for pkg name %s and version %s", pkgName, pkgVersion))
	// add dependency to cache
	d.lru.Add(fmt.Sprintf("%s:%s", pkgName, pkgVersion), npmDep)
	return npmDep, nil
}

func (d Dependencies) fetchFromRegistry(pkgName string, pkgVersion string) (*model.NpmDependency, error) {
	cSucceed := metrics.NewCounter()
	cFailure := metrics.NewCounter()
	resp, err := http.Get(fmt.Sprintf("%s/%s/%s", NpmRegistry, pkgName, pkgVersion))
	var npmDep model.NpmDependency
	if err != nil {
		err := metrics.Register(fmt.Sprintf("fetch.package.registry.%s:%s.failure", pkgName, pkgVersion), cFailure)
		if err != nil {
			d.log.Error(fmt.Sprintf("fail to log fetch.package.registry.%s:%s.failure metric", pkgName, pkgVersion))
		}
		cFailure.Inc(47)
		return nil, fmt.Errorf("failed to fetch pakcge data from npm registry: %s", err.Error())
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			d.log.Error("failed to close input steam")
		}
	}()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to fetch pakcge data from npm registry")
	}
	err = json.NewDecoder(resp.Body).Decode(&npmDep)
	if err != nil {
		err := metrics.Register(fmt.Sprintf("fetch.package.registry.%s:%s.failure", pkgName, pkgVersion), cFailure)
		if err != nil {
			d.log.Error(fmt.Sprintf("fail to log fetch.package.registry.%s:%s.failure metric", pkgName, pkgVersion))
		}
		cFailure.Inc(47)
		return nil, fmt.Errorf("failed to decode pakcge data: %s", err.Error())
	}
	err = metrics.Register(fmt.Sprintf("fetch.package.registry.%s:%s.succeeded", pkgName, pkgVersion), cSucceed)
	if err != nil {
		d.log.Error(fmt.Sprintf("fail to log fetch.package.registry.%s:%s.succeeded metric", pkgName, pkgVersion))
	}
	cSucceed.Inc(47)
	return &npmDep, nil
}

func (d Dependencies) getLru() *cache.Lru {
	return d.lru
}
