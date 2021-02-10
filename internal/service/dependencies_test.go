package service

import (
	"github.com/chen-keinan/npm-dep-tree/internal/cache"
	"github.com/chen-keinan/npm-dep-tree/internal/logger"
	"github.com/chen-keinan/npm-dep-tree/pkg/model"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var d Dep

// TestMain will exec each test, one by one
func TestMain(m *testing.M) {
	// exec setUp function
	setUp()
	// exec test and this returns an exit code to pass to os
	retCode := m.Run()
	os.Exit(retCode)
}

func setUp() {
	d = NewDependencies(logger.NewZapLogger(), cache.NewLru())
}

func Test_ResolveDependenciesAddToCache(t *testing.T) {
	dt := &model.DependencyTree{Name: "async", Version: "2.0.1", Dependencies: []*model.DependencyTree{}}
	err := d.ResolveDependencies(dt)
	assert.NoError(t, err)
	assert.Equal(t, d.getLru().Len(), 2)
	assert.Equal(t, dt.Name, "async")
	assert.Equal(t, dt.Version, "2.0.1")
	assert.Equal(t, len(dt.Dependencies), 1)
	assert.Equal(t, dt.Dependencies[0].Name, "lodash")
}

func Test_ResolveDependenciesFromCache(t *testing.T) {
	dt := &model.DependencyTree{Name: "async", Version: "2.0.1", Dependencies: []*model.DependencyTree{}}
	err := d.ResolveDependencies(dt)
	assert.NoError(t, err)
	assert.Equal(t, d.getLru().Len(), 2)
	assert.Equal(t, dt.Name, "async")
	assert.Equal(t, dt.Version, "2.0.1")
	assert.Equal(t, len(dt.Dependencies), 1)
	assert.Equal(t, dt.Dependencies[0].Name, "lodash")
}
