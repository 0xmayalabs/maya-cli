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

func TestBenchmarkRotate270(t *testing.T) {
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
			name:           "rotate270_xsmall",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       10,
			heightNew:      10,
			backend:        "groth16",
		},
		{
			name:           "rotate270_small",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       100,
			heightNew:      100,
			backend:        "groth16",
		},
		{
			name:           "rotate270_medium",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       250,
			heightNew:      250,
			backend:        "groth16",
		},
		{
			name:           "rotate270_large",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       500,
			heightNew:      500,
			backend:        "groth16",
		},
		{
			name:           "rotate270_xlarge",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       750,
			heightNew:      750,
			backend:        "groth16",
		},
		{
			name:           "rotate270_xsmall",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       10,
			heightNew:      10,
			backend:        "plonk",
		},
		{
			name:           "rotate270_small",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       100,
			heightNew:      100,
			backend:        "plonk",
		},
		{
			name:           "rotate270_medium",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       250,
			heightNew:      250,
			backend:        "plonk",
		},
		{
			name:           "rotate270_large",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       500,
			heightNew:      500,
			backend:        "plonk",
		},
		{
			name:           "rotate270_xlarge",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       750,
			heightNew:      750,
			backend:        "plonk",
		},
	}

	mdFilePath := path.Join(*resultsDir, "rotate270.md")
	mdFile, err := os.Create(mdFilePath)
	require.NoError(t, err)

	fmt.Fprintln(mdFile, "## Rotate 270")
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
			rotate270Image(t, croppedImg, finalImg)

			conf := rotate270Config{
				originalImg:  croppedImg,
				finalImg:     finalImg,
				proofDir:     dir,
				markdownFile: mdFilePath,
				backend:      tt.backend,
			}
			err := proveRotate270(conf)
			require.NoError(t, err)
		})
	}
}

func rotate270Image(t *testing.T, original, final string) {
	t0 := time.Now()
	imgFile, err := os.Open(original)
	require.NoError(t, err)
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	require.NoError(t, err)

	bounds := img.Bounds()
	rotated := image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx())) // New image with swapped dimensions

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dstX := y
			dstY := bounds.Dx() - 1 - x
			rotated.Set(dstX, dstY, img.At(x, y))
		}
	}

	// Create a new file to save the cropped image
	outFile, err := os.Create(final)
	require.NoError(t, err)
	defer outFile.Close()

	err = png.Encode(outFile, rotated)
	require.NoError(t, err)

	fmt.Printf("Time taken to rotate270: %v\n", time.Now().Sub(t0))
}

func TestRotate270(t *testing.T) {
	proofDir := t.TempDir()
	conf := rotate270Config{
		originalImg: "../sample/original.png",
		finalImg:    "../sample/rotated270.png",
		proofDir:    proofDir,
	}

	err := proveRotate270(conf)
	require.NoError(t, err)

	verifyConf := verifyRotate270Config{
		finalImg: "../sample/rotated270.png",
		proofDir: proofDir,
	}

	err = verifyRotate270(verifyConf)
	require.NoError(t, err)
}
