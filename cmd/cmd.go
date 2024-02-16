package cmd

import (
	"github.com/spf13/cobra"
)

// New returns a new cobra command that handles stackr aggregator commands and subcommands.
func New() *cobra.Command {
	return newRootCmd(
		newProveCmd(newCropCmd()),
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
