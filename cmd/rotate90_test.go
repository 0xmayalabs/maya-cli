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

func TestRotate90Benchmark(t *testing.T) {
	tests := []struct {
		name           string
		originalImg    string
		widthStartNew  int
		heightStartNew int
		widthNew       int
		heightNew      int
	}{
		{
			name:           "rotate90_xsmall",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       10,
			heightNew:      10,
		},
		{
			name:           "rotate90_small",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       100,
			heightNew:      100,
		},
		{
			name:           "rotate90_medium",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       250,
			heightNew:      250,
		},
		{
			name:           "rotate90_large",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       500,
			heightNew:      500,
		},
		{
			name:           "rotate90_xlarge",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       750,
			heightNew:      750,
		},
	}

	mdFilePath := "../book/rotate90.md"
	mdFile, err := os.Create(mdFilePath)
	require.NoError(t, err)

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
			rotate90Image(t, croppedImg, finalImg)

			conf := rotate90Config{
				originalImg:  croppedImg,
				finalImg:     finalImg,
				proofDir:     dir,
				markdownFile: mdFilePath,
			}
			err := proveRotate90(conf)
			require.NoError(t, err)
		})
	}
}

func rotate90Image(t *testing.T, original, final string) {
	t0 := time.Now()
	imgFile, err := os.Open(original)
	require.NoError(t, err)
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	require.NoError(t, err)

	bounds := img.Bounds()
	rotated := image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			srcX := x - bounds.Min.X
			srcY := y - bounds.Min.Y
			rotated.Set(bounds.Max.Y-srcY-1, srcX, img.At(x, y))
		}
	}

	// Create a new file to save the cropped image
	outFile, err := os.Create(final)
	require.NoError(t, err)
	defer outFile.Close()

	err = png.Encode(outFile, rotated)
	require.NoError(t, err)

	fmt.Printf("Time taken to rotate90: %v\n", time.Now().Sub(t0))
}

func TestRotate90(t *testing.T) {
	proofDir := t.TempDir()
	conf := rotate90Config{
		originalImg: "../sample/original.png",
		finalImg:    "../sample/rotated90.png",
		proofDir:    proofDir,
	}

	err := proveRotate90(conf)
	require.NoError(t, err)

	verifyConf := verifyRotate90Config{
		finalImg: "../sample/rotated90.png",
		proofDir: proofDir,
	}

	err = verifyRotate90(verifyConf)
	require.NoError(t, err)
}
