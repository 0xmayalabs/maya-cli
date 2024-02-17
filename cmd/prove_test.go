package cmd

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestProveCrop(t *testing.T) {
	conf := cropConfig{
		originalImg:    "original.png",
		croppedImg:     "cropped.png",
		widthStartNew:  0,
		heightStartNew: 0,
		proofDir:       ".",
	}

	err := proveCrop(conf)
	require.NoError(t, err)
}

func TestImgToPixel(t *testing.T) {
	originalImage, err := os.Open("original.png")
	require.NoError(t, err)
	defer originalImage.Close()

	_, err = convertImgToPixels(originalImage)
	require.NoError(t, err)
}
