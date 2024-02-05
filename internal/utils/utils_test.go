package utils_test

import (
	"testing"

	"github.com/anjolaoluwaakindipe/duller/internal/utils"
	"github.com/stretchr/testify/assert"
)

func Test_MakePathValid(t *testing.T) {
	t.Run("WHEN given valid string SHOULD string should remain the same", func(t *testing.T) {
		t.Parallel()
		input := "/hello"

		utils.MakeUrlPathValid(&input)

		assert.Equal(t, "/hello", input)
	})

	t.Run("WHEN given invalid string withough '/' as first character", func(t *testing.T) {
		t.Parallel()
		input := "hello/"

		utils.MakeUrlPathValid(&input)

		assert.Equal(t, "/hello", input)
	})
}
