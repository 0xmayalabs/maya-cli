package cmd

import (
	"github.com/stretchr/testify/require"
	"image"
	"image/color"
	"image/png"
	"os"
	"path"
	"testing"
)

func TestCreateBrighten(t *testing.T) {
	// t.Skip()

	inputFile, err := os.Open("../sample/original.png")
	require.NoError(t, err)

	defer inputFile.Close()

	img, err := png.Decode(inputFile)
	require.NoError(t, err)

	brightenedImg := brighten(img, 2)

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

func TestBrightenE2E(t *testing.T) {
	testDir := t.TempDir()
	createPNG(t, testDir)

	// Open input file
	inputFile, err := os.Open(path.Join(testDir, "original.png"))
	require.NoError(t, err)
	defer inputFile.Close()

	originalImg, err := png.Decode(inputFile)
	require.NoError(t, err)

	brightenFactor := 2

	// Save output file
	brightenedImg := brighten(originalImg, brightenFactor)
	outputFile, err := os.Create(path.Join(testDir, "brightened.png")) // Output file
	require.NoError(t, err)
	defer outputFile.Close()

	require.NoError(t, png.Encode(outputFile, brightenedImg))

	// Prove brighten
	conf := brightenConfig{
		originalImg:       path.Join(testDir, "original.png"),
		finalImg:          path.Join(testDir, "brightened.png"),
		brighteningFactor: 2,
		proofDir:          testDir,
	}

	err = proveBrighten(conf)
	require.NoError(t, err)
}

// brighten creates a vertically flipped version of the given image.
func brighten(img image.Image, brightnessFactor int) image.Image {
	// Create a new image for the output
	bounds := img.Bounds()
	brightenedImg := image.NewRGBA(bounds)

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

// createPNG creates a new PNG image and stores in the provided directory.
func createPNG(t *testing.T, dir string) {
	t.Helper()

	// Create a new 10x10 image
	img := image.NewGray(image.Rect(0, 0, 10, 10))

	// Iterate over the image's pixels
	for i := img.Rect.Min.Y; i < img.Rect.Max.Y; i++ {
		for j := img.Rect.Min.X; j < img.Rect.Max.X; j++ {
			// Set the pixel value to i+j
			// The value is normalized to fit within the 0-255 grayscale range
			val := uint8((i + j) * 255 / (img.Rect.Max.X + img.Rect.Max.Y - 2))
			img.SetGray(j, i, color.Gray{Y: val})
		}
	}

	// Create a new file
	f, err := os.Create(path.Join(dir, "original.png"))
	require.NoError(t, err)
	defer f.Close()

	// Encode the image to the file as PNG
	err = png.Encode(f, img)
	require.NoError(t, err)
}