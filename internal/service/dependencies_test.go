package service

import (
	"fmt"
	"github.com/chen-keinan/npm-dep-tree/internal/cache"
	"github.com/chen-keinan/npm-dep-tree/internal/logger"
	"github.com/chen-keinan/npm-dep-tree/internal/mocks"
	"github.com/chen-keinan/npm-dep-tree/internal/nhttp"
	"github.com/chen-keinan/npm-dep-tree/pkg/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
	d = NewDependencies(logger.NewZapLogger(), cache.NewLru(), nhttp.NewNpmHTTPClient())
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
	assert.Equal(t, dt.Dependencies[0].Version, "4.8.0")
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
	assert.Equal(t, dt.Dependencies[0].Version, "4.8.0")
}

func Test_ResolveDependenciesNonEscapedPkgName(t *testing.T) {
	dt := &model.DependencyTree{Name: "async/test", Version: "2.0.1", Dependencies: []*model.DependencyTree{}}
	err := d.ResolveDependencies(dt)
	assert.Error(t, err)
}

func Test_ResolveDependenciesWithMockHttp(t *testing.T) {
	responseOne := `{"name":"async","version":"2.0.1","dependencies":{"lodash":"4.8.0"}}`
	responseTwo := `{"name":"lodash","version":"4.8.0"}`
	dt := &model.DependencyTree{Name: "async", Version: "2.0.1"}
	ctl := gomock.NewController(t)
	nhttp := mocks.NewMockHClient(ctl)
	pkgNameOne := "async"
	pkgVersionOne := "2.0.1"
	pkgNameTwo := "lodash"
	pkgVersionTwo := "4.8.0"
	urlOne := fmt.Sprintf("%s/%s/%s", NpmRegistry, pkgNameOne, pkgVersionOne)
	urlTwo := fmt.Sprintf("%s/%s/%s", NpmRegistry, pkgNameTwo, pkgVersionTwo)
	nhttp.EXPECT().Get(urlOne).Return(&http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(responseOne))}, nil).Times(1)
	nhttp.EXPECT().Get(urlTwo).Return(&http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(responseTwo))}, nil).Times(1)
	// fetch data from mock registry
	lru := cache.NewLru()
	err := NewDependencies(logger.NewZapLogger(), lru, nhttp).ResolveDependencies(dt)
	assert.NoError(t, err)
	assert.Equal(t, dt.Dependencies[0].Name, "lodash")
	assert.Equal(t, dt.Dependencies[0].Version, "4.8.0")
	// fetch data from lru cache
	dtTwo := &model.DependencyTree{Name: "async", Version: "2.0.1"}
	err = NewDependencies(logger.NewZapLogger(), lru, nhttp).ResolveDependencies(dtTwo)
	assert.NoError(t, err)
	assert.Equal(t, dt.Dependencies[0].Name, "lodash")
	assert.Equal(t, dt.Dependencies[0].Version, "4.8.0")
}
