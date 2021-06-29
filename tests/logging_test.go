package tests

import (
	"github.com/HasanShahjahan/go-guest/api/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatter(t *testing.T) {
	assert.Equal(t, logging.GetFormatter("", "go"), "go")
	assert.Equal(t, logging.GetFormatter("go", "ground"), "go : ground")
	assert.Equal(t, logging.GetFormatter("go", "%v"), "go : %v")
}
