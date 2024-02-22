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

	verifyConf := verifyRotate90Config{
		finalImg: "../sample/rotated90.png",
	}

	err = verifyRotate90Crop(verifyConf)
	require.NoError(t, err)
}
