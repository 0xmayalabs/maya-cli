package cmd

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRotate180(t *testing.T) {
	proofDir := t.TempDir()
	conf := rotate180Config{
		originalImg: "../sample/original.png",
		finalImg:    "../sample/rotated180.png",
		proofDir:    proofDir,
	}

	err := proveRotate180(conf)
	require.NoError(t, err)

	verifyConf := verifyRotate180Config{
		finalImg: "../sample/rotated180.png",
		proofDir: proofDir,
	}

	err = verifyRotate180Crop(verifyConf)
	require.NoError(t, err)
}
