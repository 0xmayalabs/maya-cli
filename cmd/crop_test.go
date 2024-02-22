package cmd

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestProveCrop(t *testing.T) {
	conf := cropConfig{
		originalImg:    "../sample/original.png",
		croppedImg:     "../sample/cropped.png",
		widthStartNew:  0,
		heightStartNew: 0,
		proofDir:       t.TempDir(),
	}

	err := proveCrop(conf)
	require.NoError(t, err)
}

func TestImgToPixel(t *testing.T) {
	originalImage, err := os.Open("../sample/original.png")
	require.NoError(t, err)
	defer originalImage.Close()

	_, err = convertImgToPixels(originalImage)
	require.NoError(t, err)
}

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
