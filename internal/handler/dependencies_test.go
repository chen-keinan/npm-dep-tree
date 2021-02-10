package handler

import (
	"github.com/chen-keinan/npm-dep-tree/internal/cache"
	"github.com/chen-keinan/npm-dep-tree/internal/common"
	"github.com/chen-keinan/npm-dep-tree/internal/logger"
	"github.com/chen-keinan/npm-dep-tree/internal/service"
	"github.com/chen-keinan/npm-dep-tree/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

//Test_ResolveDependencies test resolver dependencies handler success
func Test_ResolveDependencies_OK(t *testing.T) {
	req, err := http.NewRequest("GET", path.Join(common.API, "package-dependencies/async/2.0.1"), nil)
	assert.NoError(t, err)
	zapLogger := logger.NewZapLogger()
	sh := NewDependenciesHandler(zapLogger, service.NewDependencies(zapLogger, cache.NewLru()))
	presolveRes, err := test.InvokeRequestWithResponse(req, sh.ResolveDependencies, path.Join(common.API, "package-dependencies/{name}/{version}"))
	assert.True(t, presolveRes.Code == http.StatusOK)
}

//Test_ResolveDependencies_BadRequest test resolver dependencies handler bad request
func Test_ResolveDependencies_BadRequest(t *testing.T) {
	req, err := http.NewRequest("GET", path.Join(common.API, "package-dependencies/ / "), nil)
	assert.NoError(t, err)
	zapLogger := logger.NewZapLogger()
	sh := NewDependenciesHandler(zapLogger, service.NewDependencies(zapLogger, cache.NewLru()))
	presolveRes, err := test.InvokeRequestWithResponse(req, sh.ResolveDependencies, path.Join(common.API, "package-dependencies/{name}/{version}"))
	assert.True(t, presolveRes.Code == http.StatusBadRequest)
}
