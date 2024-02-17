package cmd

import (
	"context"
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

	err := proveCrop(context.Background(), conf)
	require.NoError(t, err)

	verifyConf := verifyCropConfig{
		originalImg: "../sample/original.png",
		croppedImg:  "../sample/cropped.png",
		proofDir:    proofDir,
	}

	err = verifyCrop(context.Background(), verifyConf)
	require.NoError(t, err)
}
