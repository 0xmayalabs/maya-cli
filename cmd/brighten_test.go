package cmd

import (
	"github.com/stretchr/testify/require"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// brighten creates a vertically flipped version of the given image.
func brighten(img image.Image) image.Image {
	// Create a new image for the output
	bounds := img.Bounds()
	brightenedImg := image.NewRGBA(bounds)
	brightnessFactor := 2

	// Iterate over each pixel to adjust brightness
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			r, g, b, a := originalColor.RGBA()

			// Increase the RGB values by the brightness amount
			// Note: RGBA() returns color components in the range [0, 65535]
			r = min((r>>8)+uint32(brightnessFactor), 255)
			g = min((g>>8)+uint32(brightnessFactor), 255)
			b = min((b>>8)+uint32(brightnessFactor), 255)

			// Set the new color to the pixel
			brightenedImg.Set(x, y, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
		}
	}

	return brightenedImg
}

func TestCreateBrighten(t *testing.T) {
	// t.Skip()

	inputFile, err := os.Open("../sample/original.png")
	require.NoError(t, err)

	defer inputFile.Close()

	img, err := png.Decode(inputFile)
	require.NoError(t, err)

	brightenedImg := brighten(img)

	outputFile, err := os.Create("../sample/brightened.png") // Output file
	require.NoError(t, err)

	defer outputFile.Close()

	require.NoError(t, png.Encode(outputFile, brightenedImg))
}

func TestBrighten(t *testing.T) {
	proofDir := t.TempDir()
	conf := brightenConfig{
		originalImg:       "../sample/original.png",
		finalImg:          "../sample/brightened.png",
		brighteningFactor: 2,
		proofDir:          proofDir,
	}

	err := proveBrighten(conf)
	require.NoError(t, err)

	verifyConf := verifyBrightenConfig{
		finalImg: "../sample/brightened.png",
		proofDir: proofDir,
	}

	err = verifyBrighten(verifyConf)
	require.NoError(t, err)
}
