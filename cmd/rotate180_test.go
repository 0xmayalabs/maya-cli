package cmd

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRotate180(t *testing.T) {
	conf := rotate180Config{
		originalImg: "../sample/original.png",
		finalImg:    "../sample/rotated180.png",
	}

	err := proveRotate180(conf)
	require.NoError(t, err)

	verifyConf := verifyRotate180Config{
		finalImg: "../sample/rotated180.png",
	}

	err = verifyRotate180Crop(verifyConf)
	require.NoError(t, err)
}
