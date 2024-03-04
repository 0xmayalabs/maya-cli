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

func TestBenchmarkFlipVertical(t *testing.T) {
	tests := []struct {
		name           string
		originalImg    string
		widthStartNew  int
		heightStartNew int
		widthNew       int
		heightNew      int
		backend        string
	}{
		{
			name:           "flip_vertical_xsmall",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       10,
			heightNew:      10,
			backend:        "groth16",
		},
		{
			name:           "flip_vertical_small",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       100,
			heightNew:      100,
			backend:        "groth16",
		},
		{
			name:           "flip_vertical_medium",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       250,
			heightNew:      250,
			backend:        "groth16",
		},
		{
			name:           "flip_vertical_large",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       500,
			heightNew:      500,
			backend:        "groth16",
		},
		{
			name:           "flip_vertical_xlarge",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       750,
			heightNew:      750,
			backend:        "groth16",
		},
		{
			name:           "flip_vertical_xsmall",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       10,
			heightNew:      10,
			backend:        "plonk",
		},
		{
			name:           "flip_vertical_small",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       100,
			heightNew:      100,
			backend:        "plonk",
		},
		{
			name:           "flip_vertical_medium",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       250,
			heightNew:      250,
			backend:        "plonk",
		},
		{
			name:           "flip_vertical_large",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       500,
			heightNew:      500,
			backend:        "plonk",
		},
		{
			name:           "flip_vertical_xlarge",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       750,
			heightNew:      750,
			backend:        "plonk",
		},
	}

	mdFilePath := path.Join(*resultsDir, "flip-vertical.md")
	mdFile, err := os.Create(mdFilePath)
	require.NoError(t, err)

	fmt.Fprintln(mdFile, "## Flip vertical")
	// Write the Markdown table headers
	fmt.Fprintln(mdFile, "| Original Size | Circuit compilation (s) | Proving time (s) | Proof size (bytes) | Backend |")
	fmt.Fprintln(mdFile, "|---|---|---|---|---|")
	mdFile.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			croppedImg := path.Join(dir, "cropped.png")
			cropImage(t, tt.originalImg, croppedImg, tt.widthNew, tt.heightNew, 0, 0)

			finalImg := path.Join(dir, "final.png")
			flipVertical(t, croppedImg, finalImg)

			conf := flipVerticalConfig{
				originalImg:  croppedImg,
				finalImg:     finalImg,
				proofDir:     dir,
				markdownFile: mdFilePath,
				backend:      tt.backend,
			}
			err := proveFlipVertical(conf)
			require.NoError(t, err)
		})
	}
}

// flipVertical creates a vertically flipped version of the given image.
func flipVertical(t *testing.T, original, final string) {
	t.Helper()

	t0 := time.Now()
	imgFile, err := os.Open(original)
	require.NoError(t, err)
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	require.NoError(t, err)

	bounds := img.Bounds()
	flipped := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			flipped.Set(x, bounds.Max.Y-y-1, img.At(x, y))
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

func TestFlipVertical(t *testing.T) {
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
