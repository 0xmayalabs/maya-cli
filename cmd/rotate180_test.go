package cmd

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"image"
	"image/png"
	"os"
	"path"
	"testing"
	"time"
)

func TestBenchmarkRotate180(t *testing.T) {
	tests := []struct {
		name           string
		originalImg    string
		widthStartNew  int
		heightStartNew int
		widthNew       int
		heightNew      int
	}{
		{
			name:           "rotate180_xsmall",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       10,
			heightNew:      10,
		},
		{
			name:           "rotate180_small",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       100,
			heightNew:      100,
		},
		{
			name:           "rotate180_medium",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       250,
			heightNew:      250,
		},
		{
			name:           "rotate180_large",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       500,
			heightNew:      500,
		},
		{
			name:           "rotate180_xlarge",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       750,
			heightNew:      750,
		},
	}

	mdFilePath := path.Join(*resultsDir, "rotate180.md")
	mdFile, err := os.Create(mdFilePath)
	require.NoError(t, err)

	fmt.Fprintln(mdFile, "## Rotate 180")
	// Write the Markdown table headers
	fmt.Fprintln(mdFile, "| Original Size | Circuit compilation (s) | Proving time (s) | Proof size (bytes) |")
	fmt.Fprintln(mdFile, "|---|---|---|---|")
	mdFile.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			croppedImg := path.Join(dir, "cropped.png")
			cropImage(t, tt.originalImg, croppedImg, tt.widthNew, tt.heightNew)

			finalImg := path.Join(dir, "final.png")
			rotate180Image(t, croppedImg, finalImg)

			conf := rotate180Config{
				originalImg:  croppedImg,
				finalImg:     finalImg,
				proofDir:     dir,
				markdownFile: mdFilePath,
			}
			err := proveRotate180(conf)
			require.NoError(t, err)
		})
	}
}

func rotate180Image(t *testing.T, original, final string) {
	t0 := time.Now()
	imgFile, err := os.Open(original)
	require.NoError(t, err)
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	require.NoError(t, err)

	bounds := img.Bounds()
	rotated := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Calculate the new position for the pixel
			newX := bounds.Max.X - x - 1
			newY := bounds.Max.Y - y - 1

			// Set the pixel at the new position to the value of the current pixel
			rotated.Set(newX, newY, img.At(x, y))
		}
	}

	// Create a new file to save the cropped image
	outFile, err := os.Create(final)
	require.NoError(t, err)
	defer outFile.Close()

	err = png.Encode(outFile, rotated)
	require.NoError(t, err)

	fmt.Printf("Time taken to rotate180: %v\n", time.Now().Sub(t0))
}

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
