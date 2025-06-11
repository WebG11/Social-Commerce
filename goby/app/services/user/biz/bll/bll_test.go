package bll

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	assert.Equal(t, 4, Add(1, 3), "The result should be 4")
	assert.NotEqual(t, 5, Add(1, 3), "The result should not be 5")
	assert.True(t, Add(2, 2) == 4, "2 + 2 should be 4")
	assert.False(t, Add(2, 2) == 5, "2 + 2 should not be 5")
}
