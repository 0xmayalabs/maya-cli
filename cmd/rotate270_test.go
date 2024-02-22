package cmd

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRotate270(t *testing.T) {
	proofDir := t.TempDir()
	conf := rotate270Config{
		originalImg: "../sample/original.png",
		finalImg:    "../sample/rotated270.png",
		proofDir:    proofDir,
	}

	err := proveRotate270(conf)
	require.NoError(t, err)

	verifyConf := verifyRotate270Config{
		finalImg: "../sample/rotated270.png",
		proofDir: proofDir,
	}

	err = verifyRotate270Crop(verifyConf)
	require.NoError(t, err)
}
