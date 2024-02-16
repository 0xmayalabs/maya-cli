package cmd

import "github.com/spf13/cobra"

func newProveCmd(cmds ...*cobra.Command) *cobra.Command {
	root := &cobra.Command{
		Use:   "prove",
		Short: "Generates proof for the specified edit.",
		Long:  "Generates zero knowledge proof of edit on the original image resulting in a new image.",
	}

	root.AddCommand(cmds...)

	return root
}
