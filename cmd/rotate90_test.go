package cmd

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRotate90(t *testing.T) {
	conf := rotate90Config{
		originalImg: "../sample/original.png",
		finalImg:    "../sample/rotated90.png",
	}

	err := proveRotate90(conf)
	require.NoError(t, err)
}
