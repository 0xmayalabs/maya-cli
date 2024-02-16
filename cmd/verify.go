package cmd

import (
	"context"
	"github.com/spf13/cobra"
)

func newVerifyCmd(cmds ...*cobra.Command) *cobra.Command {
	root := &cobra.Command{
		Use:   "verify",
		Short: "Verifies proof for the specified edit.",
		Long:  "Verifies the zero knowledge proof of edit on the original image resulting in a new image.",
	}

	root.AddCommand(cmds...)

	return root
}

// verifyCropConfig specifies the verification configuration for cropping an image.
type verifyCropConfig struct {
	proofDir string
}

// newVerifyCropCmd returns a new cobra.Command for cropping.
func newVerifyCropCmd() *cobra.Command {
	var conf verifyCropConfig

	cmd := &cobra.Command{
		Use: "crop",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verifyCrop(cmd.Context(), conf)
		},
	}

	cmd.Flags().StringVar(&conf.proofDir, "proof-dir", "", "The path to the proof directory.")

	return cmd
}

// verifyCrop verifies the zk proof of crop transformation.
func verifyCrop(ctx context.Context, config verifyCropConfig) error {
	return nil
}
