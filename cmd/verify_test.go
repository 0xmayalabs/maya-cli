package cmd

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVerifyCrop(t *testing.T) {
	proofDir := t.TempDir()
	conf := cropConfig{
		originalImg:    "../sample/original.png",
		croppedImg:     "../sample/cropped.png",
		widthStartNew:  0,
		heightStartNew: 0,
		proofDir:       proofDir,
	}

	err := proveCrop(conf)
	require.NoError(t, err)

	verifyConf := verifyCropConfig{
		croppedImg: "../sample/cropped.png",
		proofDir:   proofDir,
	}

	err = verifyCrop(verifyConf)
	require.NoError(t, err)
}
