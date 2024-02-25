package cmd

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContrast(t *testing.T) {
	proofDir := t.TempDir()
	conf := contrastConfig{
		originalImg:    "../sample/original.png",
		finalImg:       "../sample/Contrasted.png",
		contrastFactor: 5,
		proofDir:       proofDir,
	}

	err := proveContrast(conf)
	require.NoError(t, err)

	verifyConf := verifyContrastConfig{
		finalImg: "../sample/Contrasted.png",
		proofDir: proofDir,
	}

	err = verifyContrast(verifyConf)
	require.NoError(t, err)
}
