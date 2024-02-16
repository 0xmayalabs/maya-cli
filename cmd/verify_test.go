package cmd

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVerifyCrop(t *testing.T) {
	conf := verifyCropConfig{}
	err := verifyCrop(context.Background(), conf)
	require.NoError(t, err)
}
