package cmd

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProveFlipHorizontal(t *testing.T) {
	proofDir := t.TempDir()
	conf := flipHorizontalConfig{
		originalImg: "../sample/original.png",
		finalImg:    "../sample/flipped_horizontal.png",
		proofDir:    proofDir,
	}

	err := proveFlipHorizontal(conf)
	require.NoError(t, err)

	verifyConf := verifyFlipHorizontalConfig{
		finalImg: "../sample/flipped_horizontal.png",
		proofDir: proofDir,
	}

	err = verifyFlipHorizontal(verifyConf)
	require.NoError(t, err)
}
