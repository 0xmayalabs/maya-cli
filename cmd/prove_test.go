package cmd

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProveCrop(t *testing.T) {
	conf := cropConfig{
		originalImg:    "original.png",
		croppedImg:     "cropped.png",
		widthStartNew:  0,
		heightStartNew: 0,
	}

	err := proveCrop(context.Background(), conf)
	require.NoError(t, err)
}
