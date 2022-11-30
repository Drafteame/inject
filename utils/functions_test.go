package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhereErrorOut(t *testing.T) {
	t.Run("get error index from single return value function", func(t *testing.T) {
		fun := func() error { return nil }

		index, exists := WhereErrorOut(reflect.TypeOf(fun))

		assert.True(t, exists)
		assert.Equal(t, 0, index)
	})

	t.Run("get error index from multi-return value function", func(t *testing.T) {
		fun := func() (string, error) { return "", nil }

		index, exists := WhereErrorOut(reflect.TypeOf(fun))

		assert.True(t, exists)
		assert.Equal(t, 1, index)
	})

	t.Run("get error index from multi-error return value function", func(t *testing.T) {
		fun := func() (string, error, error) { return "", nil, nil }

		index, exists := WhereErrorOut(reflect.TypeOf(fun))

		assert.True(t, exists)
		assert.Equal(t, 1, index)
	})

	t.Run("get error index from non function value", func(t *testing.T) {
		index, exists := WhereErrorOut(reflect.TypeOf("fun"))

		assert.False(t, exists)
		assert.Equal(t, 0, index)
	})

	t.Run("get error index from function that doesn't return error", func(t *testing.T) {
		fun := func() string { return "" }

		index, exists := WhereErrorOut(reflect.TypeOf(fun))

		assert.False(t, exists)
		assert.Equal(t, 0, index)
	})

	t.Run("get error index from function that doesn't return anything", func(t *testing.T) {
		fun := func() {}

		index, exists := WhereErrorOut(reflect.TypeOf(fun))

		assert.False(t, exists)
		assert.Equal(t, 0, index)
	})
}
