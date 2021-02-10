package handler

import (
	"github.com/chen-keinan/npm-dep-tree/internal/common"
	"github.com/chen-keinan/npm-dep-tree/internal/logger"
	"github.com/chen-keinan/npm-dep-tree/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

//Test_Ping test ping handler
func Test_Ping(t *testing.T) {
	req, err := http.NewRequest("GET", path.Join(common.API, "ping"), nil)
	assert.NoError(t, err)
	sh := NewSystemHandler(logger.NewZapLogger())
	pingRes, err := test.InvokeRequestWithResponse(req, sh.Ping, path.Join(common.API, "ping"))
	assert.True(t, pingRes.Code == 200)
}
