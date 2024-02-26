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

func TestFlipHorizontalBenchmark(t *testing.T) {
	tests := []struct {
		name           string
		originalImg    string
		widthStartNew  int
		heightStartNew int
		widthNew       int
		heightNew      int
	}{
		{
			name:           "flip_horizontal_xsmall",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       10,
			heightNew:      10,
		},
		{
			name:           "flip_horizontal_small",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       100,
			heightNew:      100,
		},
		{
			name:           "flip_horizontal_medium",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       250,
			heightNew:      250,
		},
		{
			name:           "flip_horizontal_large",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       500,
			heightNew:      500,
		},
		{
			name:           "flip_horizontal_xlarge",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       750,
			heightNew:      750,
		},
	}

	mdFilePath := "../book/perf/flip-horizontal.md"
	mdFile, err := os.Create(mdFilePath)
	require.NoError(t, err)

	fmt.Fprintln(mdFile, "## Flip horizontal")
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
			flipHorizontal(t, croppedImg, finalImg)

			conf := flipHorizontalConfig{
				originalImg:  croppedImg,
				finalImg:     finalImg,
				proofDir:     dir,
				markdownFile: mdFilePath,
			}
			err = proveFlipHorizontal(conf)
			require.NoError(t, err)
		})
	}
}

// flipHorizontal creates a horizontally flipped version of the given image.
func flipHorizontal(t *testing.T, original, final string) {
	t.Helper()

	t0 := time.Now()
	imgFile, err := os.Open(original)
	require.NoError(t, err)
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	require.NoError(t, err)

	bounds := img.Bounds()
	flipped := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			srcColor := img.At(x, y)
			flipped.Set(bounds.Max.X-x-1, y, srcColor) // Flip horizontally
		}
	}

	// Create a new file to save the cropped image
	outFile, err := os.Create(final)
	require.NoError(t, err)
	defer outFile.Close()

	err = png.Encode(outFile, flipped)
	require.NoError(t, err)

	fmt.Printf("Time taken to flip vertical: %v\n", time.Now().Sub(t0))
}

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
