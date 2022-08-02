package evictor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUnevictableError(t *testing.T) {
	t.Parallel()

	assert.True(t, IsUnevictableError(ErrNoDisruptions))
	assert.True(t, IsUnevictableError(ErrTooManyPDBs))
	assert.False(t, IsUnevictableError(nil))
	assert.False(t, IsUnevictableError(ErrNotFound))
	assert.False(t, IsUnevictableError(errors.New("foo")))
}
