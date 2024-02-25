package cmd

import (
	"github.com/spf13/cobra"
)

// New returns a new cobra command that handles stackr aggregator commands and subcommands.
func New() *cobra.Command {
	return newRootCmd(
		newProveCmd(
			newCropCmd(),
			newRotate90Cmd(),
			newRotate180Cmd(),
			newRotate270Cmd(),
			newFlipVerticalCmd(),
			newFlipHorizontalCmd(),
			newBrightenCmd(),
		),
		newVerifyCmd(
			newVerifyCropCmd(),
			newVerifyRotate90Cmd(),
			newVerifyRotate180Cmd(),
			newVerifyRotate270Cmd(),
			newVerifyFlipVerticalCmd(),
			newVerifyFlipHorizontalCmd(),
			newVerifyBrightenCmd(),
		),
	)
}

func newRootCmd(cmds ...*cobra.Command) *cobra.Command {
	root := &cobra.Command{
		Use:   "maya",
		Short: "Maya - Client for Maya network",
		Long:  "Maya separates fake from the real.",
	}

	root.AddCommand(cmds...)

	return root
}

func newProveCmd(cmds ...*cobra.Command) *cobra.Command {
	root := &cobra.Command{
		Use:   "prove",
		Short: "Generates proof for the specified edit.",
		Long:  "Generates zero knowledge proof of edit on the original image resulting in a new image.",
	}

	root.AddCommand(cmds...)

	return root
}

func newVerifyCmd(cmds ...*cobra.Command) *cobra.Command {
	root := &cobra.Command{
		Use:   "verify",
		Short: "Verifies proof for the specified edit.",
		Long:  "Verifies the zero knowledge proof of edit on the original image resulting in a new image.",
	}

	root.AddCommand(cmds...)

	return root
}
