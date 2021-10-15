package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takemo101/dc-scheduler/boot"
	"go.uber.org/fx"
)

func Test_Domain(t *testing.T) {
	boot.Testing(
		t,
		boot.TestOptions{
			ConfigPath: "../config.testing.yml",
			FXOption:   fx.Options(),
		},
		func() {
			t.Run("test", func(t *testing.T) {
				assert.Equal(t, nil, nil)
				assert.Equal(t, 1, 1)
			})
		},
	)
}
