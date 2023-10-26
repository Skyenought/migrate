package mutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllGoFiles(t *testing.T) {
	files := GetAllGoFiles("./")
	assert.Equal(t, len(files), 3)
}
