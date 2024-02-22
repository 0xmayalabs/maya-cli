package cmd

import (
	"github.com/stretchr/testify/require"
	"image"
	"image/png"
	"os"
	"testing"
)

// flipVertical creates a vertically flipped version of the given image.
func flipVertical(img image.Image) image.Image {
	bounds := img.Bounds()
	flipped := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			flipped.Set(x, bounds.Max.Y-y-1, img.At(x, y))
		}
	}
	return flipped
}

func TestCreateFlipVertical(t *testing.T) {
	t.Skip()

	inputFile, err := os.Open("../sample/original.png") // Replace 'input.jpg' with your image file
	require.NoError(t, err)

	defer inputFile.Close()

	img, err := png.Decode(inputFile)
	require.NoError(t, err)

	flippedImg := flipVertical(img)

	outputFile, err := os.Create("flipped_vertical.png") // Output file
	require.NoError(t, err)

	defer outputFile.Close()

	require.NoError(t, png.Encode(outputFile, flippedImg))
}

func TestProveFlipVertical(t *testing.T) {
	proofDir := t.TempDir()
	conf := flipVerticalConfig{
		originalImg: "../sample/original.png",
		finalImg:    "../sample/flipped_vertical.png",
		proofDir:    proofDir,
	}

	err := proveFlipVertical(conf)
	require.NoError(t, err)

	verifyConf := verifyFlipVerticalConfig{
		finalImg: "../sample/flipped_vertical.png",
		proofDir: proofDir,
	}

	err = verifyFlipVertical(verifyConf)
	require.NoError(t, err)
}
