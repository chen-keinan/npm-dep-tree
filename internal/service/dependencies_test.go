package service

import (
	"github.com/chen-keinan/npm-dep-tree/internal/cache"
	"github.com/chen-keinan/npm-dep-tree/internal/logger"
	"github.com/chen-keinan/npm-dep-tree/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ResolveDependenciesAddToCache(t *testing.T) {
	d := NewDependencies(logger.NewZapLogger(), cache.NewLru())
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
	d := NewDependencies(logger.NewZapLogger(), cache.NewLru())
	dt := &model.DependencyTree{Name: "async", Version: "2.0.1", Dependencies: []*model.DependencyTree{}}
	err := d.ResolveDependencies(dt)
	assert.NoError(t, err)
	assert.Equal(t, d.getLru().Len(), 2)
	assert.Equal(t, dt.Name, "async")
	assert.Equal(t, dt.Version, "2.0.1")
	assert.Equal(t, len(dt.Dependencies), 1)
	assert.Equal(t, dt.Dependencies[0].Name, "lodash")
}
