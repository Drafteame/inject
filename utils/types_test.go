package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFirstReturnType(t *testing.T) {
	t.Run("get type from function with single return", func(t *testing.T) {
		fun := func() string { return "" }

		tfun := GetFirstReturnType(fun)

		assert.NotNil(t, tfun)
		assert.Equal(t, reflect.TypeOf(""), tfun)
	})

	t.Run("get type from function with no return", func(t *testing.T) {
		fun := func() {}

		tfun := GetFirstReturnType(fun)

		assert.Nil(t, tfun)
	})

	t.Run("get type from function with multi-return values", func(t *testing.T) {
		fun := func() (string, error) { return "", nil }

		tfun := GetFirstReturnType(fun)

		assert.NotNil(t, tfun)
		assert.Equal(t, reflect.TypeOf(""), tfun)
	})

	t.Run("get type from non function", func(t *testing.T) {
		tfun := GetFirstReturnType("")

		assert.Nil(t, tfun)
	})
}

// nolint
func TestEmbedsType(t *testing.T) {
	type in struct{}
	type out struct{}

	t.Run("struct embeds specified type", func(t *testing.T) {
		type some struct {
			in
		}

		embed := EmbedsType(some{}, reflect.TypeOf(in{}))

		assert.True(t, embed)
	})

	t.Run("struct embeds specified type with reflected type as input", func(t *testing.T) {
		type some struct {
			in
		}

		embed := EmbedsType(reflect.TypeOf(some{}), reflect.TypeOf(in{}))

		assert.True(t, embed)
	})

	t.Run("struct embeds specified type with pointer as input", func(t *testing.T) {
		type some struct {
			in
		}

		embed := EmbedsType(&some{}, reflect.TypeOf(in{}))

		assert.True(t, embed)
	})

	t.Run("struct embeds specified type with multiple embedded types", func(t *testing.T) {
		type some struct {
			out
			in
		}

		embed := EmbedsType(some{}, reflect.TypeOf(in{}))

		assert.True(t, embed)
	})

	t.Run("check embed from nil input", func(t *testing.T) {
		embed := EmbedsType(nil, reflect.TypeOf(in{}))

		assert.False(t, embed)
	})

	t.Run("struct do not embed specified type", func(t *testing.T) {
		type some struct{}

		embed := EmbedsType(some{}, reflect.TypeOf(in{}))

		assert.False(t, embed)
	})
}
