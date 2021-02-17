package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEscapePackageName(t *testing.T) {
	ep := EscapePackageName("a/b")
	assert.Equal(t, ep, "a%2Fb")
	ep = EscapePackageName("a")
	assert.Equal(t, ep, "a")
}

func TestTrimVersionSign(t *testing.T) {
	ep := TrimVersionSign("~1.0.1")
	assert.Equal(t, ep, "1.0.1")
	ep = TrimVersionSign("^1.0.1")
	assert.Equal(t, ep, "1.0.1")
	ep = TrimVersionSign("1.0.1")
	assert.Equal(t, ep, "1.0.1")
}
