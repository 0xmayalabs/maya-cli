package cmd

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path"
	"testing"
	"time"
)

func TestBenchmarkCrop(t *testing.T) {
	tests := []struct {
		name           string
		originalImg    string
		widthStartNew  int
		heightStartNew int
		widthNew       int
		heightNew      int
	}{
		{
			name:           "crop_xsmall",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       10,
			heightNew:      10,
		},
		{
			name:           "crop_small",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       100,
			heightNew:      100,
		},
		{
			name:           "crop_medium",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       250,
			heightNew:      250,
		},
		{
			name:           "crop_large",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       500,
			heightNew:      500,
		},
		{
			name:           "crop_xlarge",
			originalImg:    "../sample/original-1000x1000.png",
			widthStartNew:  0,
			heightStartNew: 0,
			widthNew:       750,
			heightNew:      750,
		},
	}

	mdFilePath := path.Join(*resultsDir, "crop.md")
	mdFile, err := os.Create(mdFilePath)
	require.NoError(t, err)

	fmt.Fprintln(mdFile, "## Crop")
	// Write the Markdown table headers
	fmt.Fprintln(mdFile, "| Original Size | Final Size | Circuit compilation (s) | Proving time (s) | Proof size (bytes) |")
	fmt.Fprintln(mdFile, "|---|---|---|---|---|")
	mdFile.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			finalImg := path.Join(dir, "final.png")
			cropImage(t, tt.originalImg, finalImg, tt.widthNew, tt.heightNew)

			conf := cropConfig{
				originalImg:    tt.originalImg,
				croppedImg:     finalImg,
				widthStartNew:  tt.widthStartNew,
				heightStartNew: tt.heightStartNew,
				proofDir:       dir,
				markdownFile:   mdFilePath,
			}
			err := proveCrop(conf)
			require.NoError(t, err)
		})
	}
}

func cropImage(t *testing.T, original, final string, widthNew, heightNew int) {
	t0 := time.Now()
	imgFile, err := os.Open(original)
	require.NoError(t, err)
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	require.NoError(t, err)

	// The rectangle is defined by the top-left and bottom-right points: (x0, y0, x1, y1)
	cropRect := image.Rect(0, 0, widthNew, heightNew)

	// Create a new blank image with the size of the crop rectangle
	croppedImg := image.NewRGBA(cropRect)

	// Crop the image by drawing it on the new blank image
	draw.Draw(croppedImg, croppedImg.Bounds(), img, cropRect.Min, draw.Src)

	// Create a new file to save the cropped image
	outFile, err := os.Create(final)
	require.NoError(t, err)
	defer outFile.Close()

	err = png.Encode(outFile, croppedImg)
	require.NoError(t, err)

	fmt.Printf("Time taken to crop: %v\n", time.Now().Sub(t0))
}

func TestCrop(t *testing.T) {
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
